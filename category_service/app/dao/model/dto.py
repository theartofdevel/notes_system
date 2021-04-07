from dataclasses import dataclass

from dao.model.base import Base


@dataclass
class CreateCategoryDTO(Base):
    name: str
    user_uuid: str
    parent_uuid: str = None


@dataclass
class UpdateCategoryDTO(Base):
    name: str
    user_uuid: str
    uuid: str = None


@dataclass
class DeleteCategoryDTO(Base):
    user_uuid: str
    uuid: str = None
