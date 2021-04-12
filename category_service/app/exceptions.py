from enum import Enum


class AppError(Enum):
    def __init__(self, code: str, message: str, developer_message: str):
        self.code = code
        self.message = message
        self.developer_message = developer_message

    SYSTEM_ERROR = ("CS-00001", "system error", "")
    CATEGORY_NOT_FOUND = ("CS-00008", "category not found", "")
    USER_NOT_FOUND = ("CS-00009", "user not found", "")
    VALIDATION_ERROR = ("CS-00010", "validation error", "")


class AppException(Exception):
    def __init__(self,
                 exc_data: AppError = None,
                 code: str = None,
                 error: str = None,
                 developer_message: str = None,
                 *args):
        if exc_data:
            self.code = exc_data.code
            self.message = exc_data.message
            self.developer_message = exc_data.developer_message

        if code:
            self.code = code
        if error:
            self.message = error
        if developer_message:
            self.developer_message = developer_message

        super().__init__(*args)


class NotFoundException(AppException):
    pass


class ValidationException(AppException):
    pass
