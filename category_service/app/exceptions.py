from enum import Enum
from http import HTTPStatus


class AppError(Enum):
    def __init__(self, error_code: str, error: str, developer_message: str, http_code: int):
        self.error_code = error_code
        self.error = error
        self.developer_message = developer_message
        self.http_code = http_code

    SYSTEM_ERROR = ("CS-00001", "System Error", "", HTTPStatus.INTERNAL_SERVER_ERROR)

    CATEGORY_NOT_FOUND = ("CS-00008", "Specified category not found. Check identifier.", "", HTTPStatus.NOT_FOUND)
    USER_NOT_FOUND = ("CS-00009", "Specified user not found. Check identifier.", "", HTTPStatus.NOT_FOUND)


class AppException(Exception):
    def __init__(self,
                 exc_data: AppError = None,
                 error_code: str = None,
                 error: str = None,
                 developer_message: str = None,
                 http_code: int = HTTPStatus.INTERNAL_SERVER_ERROR,
                 *args):
        self.http_code = http_code
        if exc_data:
            self.error_code = exc_data.error_code
            self.error = exc_data.error
            self.developer_message = exc_data.developer_message
            self.http_code = exc_data.http_code

        if error_code:
            self.error_code = error_code
        if error:
            self.error = error
        if developer_message:
            self.developer_message = developer_message

        super().__init__(*args)
