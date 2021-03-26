from dataclasses import dataclass, asdict
from typing import Dict, Any

import marshmallow_dataclass


@dataclass
class Base:
    def to_dict(self) -> Dict[str, Any]:
        return {k: v for k, v in asdict(self).items() if v is not None}

    @classmethod
    def to_schema(cls):
        return marshmallow_dataclass.class_schema(cls)()