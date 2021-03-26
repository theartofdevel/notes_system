from typing import List, Dict

from py2neo import Graph

from dao.storage.storage import Storage


class Neo4jStorage(Storage):
    __slots__ = ["_graph", "url", "_username", "_password"]

    def __init__(self, hostname: str, password: str, port: int = 7687, username: str = "neo4j"):
        self.url = f"bolt://{hostname}:{port}"
        self._username = username
        self._password = password

        self._connect()

    def _connect(self):
        self._graph = Graph(self.url, password=self._password, user=self._username)

    def _execute_cypher(self, command: str) -> List:
        cur = self._graph.run(command)
        return cur.data()

    def find_one(self, cypher_cmd: str) -> List:
        data = self._execute_cypher(command=cypher_cmd)
        return data

    def find(self, cypher_cmd: str) -> List:
        data = self._execute_cypher(command=cypher_cmd)
        return data

    def create(self, cypher_cmd: str) -> List:
        data = self._execute_cypher(command=cypher_cmd)
        return data

    def update(self, cypher_cmd: str) -> List:
        data = self._execute_cypher(command=cypher_cmd)
        return data

    def delete(self, cypher_cmd: str) -> List:
        data = self._execute_cypher(command=cypher_cmd)
        return data

    def execute(self, cypher_cmd: str) -> List:
        return self._execute_cypher(command=cypher_cmd)