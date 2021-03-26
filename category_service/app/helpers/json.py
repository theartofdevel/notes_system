import collections
import enum
import json
from typing import Any


class CustomJSONEncoder(json.JSONEncoder):
    def default(self, obj):
        return self.get_dict(obj=obj)

    def get_dict(self, obj) -> Any:
        if obj is None:
            return
        if isinstance(obj, str):
            return obj
        elif isinstance(obj, enum.Enum):
            return str(obj)
        elif isinstance(obj, dict):
            return dict((key, self.get_dict(val)) for key, val in obj.items() if val is not None)
        elif isinstance(obj, collections.Iterable):
            return [self.get_dict(val) for val in obj]
        elif hasattr(obj, '__slots__'):
            return self.get_dict(dict(
                (name, getattr(obj, name)) for name in getattr(obj, '__slots__') if getattr(obj, name) is not None))
        elif hasattr(obj, '__dict__'):
            return self.get_dict(vars(obj))
        else:
            return json.JSONEncoder.default(self, obj)
