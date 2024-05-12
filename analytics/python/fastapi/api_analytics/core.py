import os
import threading
from datetime import datetime
from typing import Dict, List
from logging import getLogger

import requests

logger = getLogger(__name__)
DEFAULT_SERVER_URL = os.getenv("ANALYTICS_SERVER_URL", "https://www.apianalytics-server.com")

_requests = []
_last_posted = datetime.now()


def _post_requests(
    api_key: str, requests_data: List[Dict], framework: str, privacy_level: int
):
    logger.info("Analytics: Sending requests to server")
    logger.info(f"Analytics-url:{DEFAULT_SERVER_URL}/api/log-request")
    res = requests.post(
        f"{DEFAULT_SERVER_URL}/api/log-request",
        json={
            "api_key": api_key,
            "requests": requests_data,
            "framework": framework,
            "privacy_level": privacy_level,
        },
        timeout=10,
    )

    logger.info("Analytics: request to server: %s", res.request.body)
    logger.info("Analytics: Response from server: %s", res.status_code)


def log_request(api_key: str, request_data: Dict, framework: str, privacy_level: int):
    if api_key == "" or api_key is None:
        return
    global _requests, _last_posted
    _requests.append(request_data)
    now = datetime.now()
    if (now - _last_posted).total_seconds() > 60.0:
        threading.Thread(
            target=_post_requests, args=(api_key, _requests, framework, privacy_level)
        ).start()
        _requests = []
        _last_posted = now
