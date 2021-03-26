from logging import Logger

from injector import Module, Binder, singleton

from config import Config
from dao.category.category import CategoryDAO
from dao.category.neo4j import Neo4jCategoryDAO
from dao.storage.neo4j import Neo4jStorage
from dao.storage.storage import Storage


class StorageModule(Module):
    def __init__(self, config: Config):
        self.config = config

    def configure(self, binder: Binder) -> None:
        neo4j_storage = Neo4jStorage(hostname=self.config.NEO4J_HOSTNAME,
                                     port=self.config.NEO4J_PORT,
                                     username=self.config.NEO4J_LOGIN,
                                     password=self.config.NEO4J_PASSWORD)
        binder.bind(interface=Storage, to=neo4j_storage, scope=singleton)
        binder.bind(interface=CategoryDAO, to=Neo4jCategoryDAO, scope=singleton)


class LoggerModule(Module):
    def __init__(self, logger: Logger):
        self.logger = logger

    def configure(self, binder: Binder) -> None:
        binder.bind(interface=Logger, to=self.logger, scope=singleton)
