from typing import List

from injector import inject

from dao.category.category import CategoryDAO
from dao.model.dto import DeleteCategoryDTO, UpdateCategoryDTO, CreateCategoryDTO
from dao.model.model import Category
from dao.storage.storage import Storage


class Neo4jCategoryDAO(CategoryDAO):
    __slots__ = ["storage"]

    @inject
    def __init__(self, storage: Storage):
        self.storage = storage

    def find_user_categories(self, user_uuid: str) -> List["Category"]:
        result = self.storage.find(
            f"""
            MATCH path = (u:User)-[*]->(c)
            WHERE NOT (c)-->() AND u.id = "{user_uuid}"
            WITH collect(path) AS ps
            CALL apoc.convert.toTree(ps) yield value
            RETURN value
            """
        )
        return self._parse_categories(categories=result[0]["value"]["own"])

    def _parse_categories(self, categories: List, parent_uuid: str = None) -> List["Category"]:
        res = []
        for el in categories:
            children = None
            if "child" in el:
                children = self._parse_categories(categories=el["child"], parent_uuid=el["id"])
            c = Category(uuid=el["id"], name=el["name"], parent_uuid=parent_uuid, children=children)
            res.append(c)
        return res

    def check_user_exist(self, user_uuid: str) -> bool:
        return self._check_entity_exist(entity="User", euuid=user_uuid)

    def check_category_exist(self, category_uuid: str) -> bool:
        return self._check_entity_exist(entity="Category", euuid=category_uuid)

    def _check_entity_exist(self, entity: str, euuid: str) -> bool:
        result = self.storage.execute(
            f"""
            OPTIONAL MATCH (n:{entity} {'{'} id: "{euuid}" {'}'})
            RETURN n IS NOT NULL AS is_exist
            """
        )
        return result[0]["is_exist"]

    def create_root_category(self, category: CreateCategoryDTO):
        result = self.storage.create(
            f"""
            MERGE (u:User {'{'} id: "{category.user_uuid}" {'}'})
            CREATE (c:Category {'{'} name: "{category.name}", id: apoc.create.uuid() {'}'})
            CREATE (u)-[r:OWN]->(c)
            RETURN c.id AS category_id
            """
        )
        return Category(uuid=result[0]["category_id"], name=category.name, parent_uuid=category.parent_uuid)

    def create_sub_category(self, category: CreateCategoryDTO):
        result = self.storage.create(
            f"""
            MATCH (cs:Category)
            WHERE cs.id = "{category.parent_uuid}"
            CREATE (c:Category {'{'} name: "{category.name}", id: apoc.create.uuid() {'}'})
            CREATE (cs)-[r:CHILD]->(c)
            RETURN c.id AS category_id
            """
        )
        return Category(uuid=result[0]["category_id"], name=category.name, parent_uuid=category.parent_uuid)

    def update_category(self, category: UpdateCategoryDTO):
        # TODO update parent id
        self.storage.update(
            f"""
            MATCH (c:Category {'{'} id: "{category.uuid}" {'}'})
            SET c.name = "{category.name}"
            """
        )

    def delete_category(self, category: DeleteCategoryDTO):
        self.storage.update(
            f"""
            MATCH path = (c:Category)-[*0..]->(cc:Category)
            WHERE c.id = "{category.uuid}"
            DETACH DELETE path
            """
        )
