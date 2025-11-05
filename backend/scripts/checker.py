#!/usr/bin/env python3
import requests
import random
import string

BASE = "http://localhost:8080"
SESSION = requests.Session()


def rand_suffix(n=5):
    return ''.join(random.choice(string.ascii_lowercase + string.digits) for _ in range(n))


def print_title(title):
    print("\n" + "=" * 60)
    print(title)
    print("=" * 60)


# ---------------------------
# basic endpoints
# ---------------------------
def try_health():
    print_title("0) HEALTH")
    try:
        r = SESSION.get(f"{BASE}/api/v1/health")
        print("status:", r.status_code)
        print("body:", r.text)
    except Exception as e:
        print("health failed:", e)


def try_register(email, username, password):
    print_title("1) REGISTER")
    url = f"{BASE}/api/v1/auth/register"
    payload = {
        "email": email,
        "username": username,
        "password": password
    }
    r = SESSION.post(url, json=payload)
    print("status:", r.status_code)
    print("body:", r.text)
    return r


def try_login(login, password):
    print_title("2) LOGIN")
    url = f"{BASE}/api/v1/auth/login"
    payload = {
        "login": login,   # у тебя логин по полю login
        "password": password,
    }
    r = SESSION.post(url, json=payload)
    print("status:", r.status_code)
    print("body:", r.text)
    return r


def try_auth_me():
    print_title("2.5) AUTH ME")
    url = f"{BASE}/api/v1/auth/me"
    r = SESSION.get(url)
    print("status:", r.status_code)
    print("body:", r.text)
    return r


def try_get_storage_settings():
    print_title("3) GET STORAGE SETTINGS")
    url = f"{BASE}/api/v1/settings/storage"
    r = SESSION.get(url)
    print("status:", r.status_code)
    print("body:", r.text)
    return r


def try_set_storage_settings(storage_name):
    print_title("4) SET STORAGE SETTINGS")
    url = f"{BASE}/api/v1/settings/storage"
    r = SESSION.post(url, json={"storage": storage_name})
    print("status:", r.status_code)
    print("body:", r.text)
    return r


# ---------------------------
# paste helpers
# ---------------------------
def try_create_paste(slug, content="print('hello')", folder="project-euler/001-100", title="first paste"):
    print_title(f"CREATE PASTE slug={slug}")
    url = f"{BASE}/api/v1/pastes"
    payload = {
        "title": title,
        "content": content,
        "extension": "py",
        "folder": folder,
        "slug": slug,
        "is_public": True,
    }
    r = SESSION.post(url, json=payload)
    print("status:", r.status_code)
    print("body:", r.text)
    return r


def try_create_anon_paste(content="print('anon')"):
    print_title("CREATE ANON PASTE")
    url = f"{BASE}/api/v1/pastes/anon"
    payload = {
        "title": "anon paste",
        "content": content,
        "extension": "txt",
        "folder": "",
        "is_public": True,
    }
    r = SESSION.post(url, json=payload)
    print("status:", r.status_code)
    print("body:", r.text)
    return r


def try_get_raw(slug):
    print_title(f"GET RAW slug={slug}")
    url = f"{BASE}/api/v1/pastes/{slug}/raw"
    r = SESSION.get(url)
    print("status:", r.status_code)
    print("body:", r.text)
    return r


def try_delete(slug):
    print_title(f"DELETE (soft) slug={slug}")
    url = f"{BASE}/api/v1/pastes/{slug}"
    r = SESSION.delete(url)
    print("status:", r.status_code)
    print("body:", r.text)
    return r


def try_my_pastes(limit=2, offset=0):
    print_title(f"MY PASTES limit={limit} offset={offset}")
    url = f"{BASE}/api/v1/me/pastes?limit={limit}&offset={offset}"
    r = SESSION.get(url)
    print("status:", r.status_code)
    print("body:", r.text)
    return r


def try_recent_pastes(limit=2, offset=0):
    print_title(f"RECENT PASTES limit={limit} offset={offset}")
    url = f"{BASE}/api/v1/pastes/recent?limit={limit}&offset={offset}"
    r = SESSION.get(url)
    print("status:", r.status_code)
    print("body:", r.text)
    return r


def try_big_paste(slug):
    print_title("CREATE BIG PASTE (expect 400 too_large)")
    url = f"{BASE}/api/v1/pastes"
    big_content = "x" * 300_000  # больше лимита 200к
    payload = {
        "title": "too big",
        "content": big_content,
        "extension": "txt",
        "folder": "",
        "slug": slug,
        "is_public": True,
    }
    r = SESSION.post(url, json=payload)
    print("status:", r.status_code)
    print("body:", r.text)
    return r


def main():
    # один и тот же юзер, которого ты уже в БД сделал админом
    email = "user_au35d@example.com"
    username = "user_au35d"
    password = "123456"

    # для уникальных паст
    base_slug = f"my-same-slug-{rand_suffix()}"

    # 0) health
    try_health()

    # 1) рега (если уже есть — будет 500/email_taken, это ок)
    try_register(email, username, password)

    # 2) логин
    login = try_login(email, password)
    if login.status_code != 200:
        print("⚠️ login failed, но продолжаем для публичных ручек")

    # 2.5) auth/me — тут надо увидеть is_admin
    try_auth_me()

    # 3) настройки — GET
    try_get_storage_settings()

    # 4) настройки — POST (если не админ → 403)
    try_set_storage_settings("local")
    #try_set_storage_settings("firebase")

    # 5) создаём 3 пасты подряд, чтобы протестить пагинацию
    s1 = base_slug + "-1"
    s2 = base_slug + "-2"
    s3 = base_slug + "-3"

    try_create_paste(s1, "print('v1')", title="p1")
    try_create_paste(s2, "print('v2')", title="p2")
    try_create_paste(s3, "print('v3')", title="p3")

    # 6) пагинация по моим пастам
    try_my_pastes(limit=2, offset=0)
    try_my_pastes(limit=2, offset=2)

    # 7) пагинация по публичным
    try_recent_pastes(limit=2, offset=0)
    try_recent_pastes(limit=2, offset=2)

    # 8) твой сценарий мягкого удаления
    slug_soft = base_slug + "-soft"
    try_create_paste(slug_soft, "print('soft v1')")
    try_create_paste(slug_soft, "print('soft v2')")   # перезапись
    try_get_raw(slug_soft)
    try_delete(slug_soft)
    try_create_paste(slug_soft, "print('soft v3')")   # после delete
    try_get_raw(slug_soft)

    # 9) анон-паста
    try_create_anon_paste("print('anon user')")

    # 10) слишком большая паста — проверка лимита
    try_big_paste(base_slug + "-big")

    print_title("DONE")


if __name__ == "__main__":
    main()
