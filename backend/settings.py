from fastapi import FastAPI

app = FastAPI()




import os
from os.path import join, dirname
from dotenv import load_dotenv

dotenv_path = join(dirname(__file__), '.env')
load_dotenv(dotenv_path)


GITHUB_WEBHOOK_CONFIG = {P
    REPO_PATH: os.environ.get("GITHUB_REPO_PATH"),
    BRANCH: os.environ.get("GITHUB_BRANCH"),
    WEBHOOK_SECRET: os.environ.get("GITHUB_WEBHOOK_SECRET"),
    PORT: os.environ.get("GITHUB_WEBHOOK_PORT"),
    ALLOWED_IPS: os.environ.get("GITHUB_ALLOWED_IPS")
}