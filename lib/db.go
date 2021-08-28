package lib

type Database struct {
  Uri string
}

func NewDB(env Env) Database {
  return Database { Uri: env.DatabaseURI }
}