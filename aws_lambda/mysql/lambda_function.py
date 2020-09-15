import json
import logging
import pymysql
import os

logger = logging.getLogger()
logger.setLevel(logging.INFO)

#rds settings
rds_host  = os.environ("RDS_HOST")
username = os.environ("USERNAME")
password = os.environ("PASSWORD")
db_name = os.environ("DB_NAME")

def lambda_handler(event, context):
    print("event: ", event)

    try:
        conn = pymysql.connect(rds_host, user=username, passwd=password, db=db_name, connect_timeout=5)
        
        with conn.cursor() as cur:
            for record in event["Records"]:
                if "new-shortlinks" in record["eventSourceARN"]:
                    insert_new_shorturl(conn, cur, record)
                    
                elif "used-shortlinks" in record["eventSourceARN"]:
                    insert_used_shorturl(conn, cur, record)
    except Exception as e:
        logger.error(e)


def insert_new_shorturl(conn, cur, record):
    print("inserting new shorturl")
    print("record: ", record)
    payload = json.loads(record['body'])
    print(payload)
    sql = "INSERT into `BITLY_URLS` (`long_url`, `short_url`, `created_at`, `created_ip`, `user_agent`) VALUES (%s, %s, %s, %s, %s);"
    cur.execute(sql, (payload["LongUrl"], payload["ShortUrlHash"], payload["CreatedAt"], payload["CreatedIP"], payload["UserAgent"]))
    conn.commit()


def insert_used_shorturl(conn, cur, record):
    print("inserting used shorturl")
    print("record: ", record)
    payload = json.loads(record['body'])
    print("payload: ", payload)
    print(str(payload["ShortUrlHash"]))
    print(str(payload["AccessedAt"]))
    sql = "INSERT into BITLY_SHORT_URL_ACCESS_FACTS (short_url, accessed_at, accessed_ip) VALUES(%s, %s, %s);"
    cur.execute(sql, (payload["ShortUrlHash"], payload["AccessedAt"], payload["AccessedIP"],))
    conn.commit()
