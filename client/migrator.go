package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kLkA/go-migrate-tool/modules"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Migrator struct {
	migrations []Migration

	config *Config
	db     *mgo.Database
}
type Migration struct {
	Name string
	Path string
}

type MigrationHistory struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
}

func New(conn *mgo.Session, config *Config) (*Migrator, error) {
	err := config.Prepare()
	if err != nil {
		return nil, err
	}

	i := &Migrator{
		db:     conn.DB(config.DatabaseName),
		config: config,
	}
	i.loadNewMigrations()

	return i, nil
}

func (m *Migrator) loadNewMigrations() {
	collection := m.db.C(m.config.TableName)

	gopath := os.Getenv("GOPATH")
	migrationFiles, err := filepath.Glob(strings.Replace(m.config.MigrationPath, "$GOPATH", gopath, 1) + "/*.json")
	if err != nil {
		modules.Log.Error(err)
	}

	migrations := []MigrationHistory{}
	err = collection.Find(nil).All(&migrations)
	if err != nil {
		modules.Log.Error(err)
	}

	for _, file := range migrationFiles {
		tmp := strings.Split(file, "/")
		name := tmp[len(tmp)-1]

		allow := true
		for _, migration := range migrations {
			if migration.Name == name {
				allow = false
				break
			}
		}
		if allow {
			m.migrations = append(m.migrations, Migration{
				Name: name,
				Path: file,
			})
		}
	}
}

func (m *Migrator) Apply() error {
	for _, migration := range m.migrations {
		var data interface{}

		content, err := ioutil.ReadFile(migration.Path)
		if err != nil {
			return migrationFailed(migration.Name, err)
		}

		err = json.Unmarshal(content, &data)
		if err != nil {
			return migrationFailed(migration.Name, err)
		}

		err = m.ProcessCommand(data, migration.Name)
		if err != nil {
			return migrationFailed(migration.Name, err)
		}
	}

	return nil
}

func (m *Migrator) ProcessCommand(data interface{}, name string) error {
	commands, ok := data.([]interface{})
	if !ok {
		command, ok := data.(map[string]interface{})
		if !ok {
			return errors.New("wrong json data format")
		}
		commands = []interface{}{command}
	}

	for _, command := range commands {
		var cmd bson.D
		for key, value := range command.(map[string]interface{}) {
			if strings.Contains(key, "cmd:") {
				key = key[4:]
				cmd = append([]bson.DocElem{{Name: key, Value: value}}, cmd...)
			} else {
				cmd = append(cmd, bson.DocElem{Name: key, Value: value})
			}
		}

		err := m.db.Run(cmd, nil)
		if err != nil {
			return err
		}
	}

	collection := m.db.C(m.config.TableName)
	err := collection.Insert(MigrationHistory{
		ID:        bson.NewObjectId(),
		Name:      name,
		CreatedAt: time.Now(),
	})
	if err != nil {
		modules.Log.Error(fmt.Sprintf("Error during insert of mongo document: %s", err.Error()))
	}

	/* save executed file to migration history in mongo */
	modules.Log.Info(fmt.Sprintf("Migration \"%s\" has been applied", name))

	return nil
}

func migrationFailed(migration string, err error) error {
	return errors.New(fmt.Sprintf("migration failed in file: %s, with error: %s", migration, err.Error()))
}
