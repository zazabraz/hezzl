package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config any

func ReadEnv(filename string, inst Config) error {
	err := godotenv.Load(filename)
	if err != nil {
		return err
	}
	err = cleanenv.ReadEnv(inst)
	if err != nil {
		return err
	}

	return nil
}
