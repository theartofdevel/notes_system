from abc import ABC, abstractmethod
from typing import List

from dao.model.dto import *


class CategoryDAO(ABC):

    @abstractmethod
    def find_user_categories(self, user_uuid: str) -> List:
        raise NotImplementedError

    @abstractmethod
    def check_user_exist(self, user_uuid: str):
        raise NotImplementedError

    @abstractmethod
    def create_root_category(self, category: CreateCategoryDTO):
        raise NotImplementedError

    @abstractmethod
    def check_category_exist(self, category_uuid: str):
        raise NotImplementedError

    @abstractmethod
    def create_sub_category(self, category: CreateCategoryDTO):
        raise NotImplementedError

    @abstractmethod
    def update_category(self, category: UpdateCategoryDTO):
        raise NotImplementedError

    @abstractmethod
    def delete_category(self, category: DeleteCategoryDTO):
        raise NotImplementedError
