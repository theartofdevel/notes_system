from dataclasses import dataclass
from typing import List

from dao.model.base import Base


@dataclass
class Category(Base):
    uuid: str
    name: str
    parent_uuid: str = None
    children: List["Category"] = None
