from abc import ABC, abstractmethod


class Storage(ABC):
    __slots__ = []

    @abstractmethod
    def find_one(self, *args, **kwargs):
        pass

    @abstractmethod
    def find(self, *args, **kwargs):
        pass

    @abstractmethod
    def create(self, *args, **kwargs):
        pass

    @abstractmethod
    def update(self, *args, **kwargs):
        pass

    @abstractmethod
    def delete(self, *args, **kwargs):
        pass

    @abstractmethod
    def execute(self, *args, **kwargs):
        pass