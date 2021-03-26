import os
from typing import Dict, Any

import yaml


class Config:
    __slots__ = ["_path",
                 "DEBUG",
                 "NEO4J_HOSTNAME",
                 "NEO4J_PORT",
                 "NEO4J_LOGIN",
                 "NEO4J_PASSWORD"
                 ]

    def __init__(self, yaml_file: str) -> None:
        self._path = yaml_file
        self._read()

    def to_dict(self) -> Dict[str, Any]:
        return {field.lower(): getattr(self, field) for field in self.__slots__}

    def _read(self) -> None:
        if not os.path.exists(self._path):
            raise AttributeError(f"config yaml doesnt exist: {self._path}")

        with open(self._path) as config_file:
            config_content = config_file.read()
            config_yaml = yaml.safe_load(config_content)

        for k, v in config_yaml.items():
            k = k.upper()
            if k in self.__slots__:
                setattr(self, k, v)
