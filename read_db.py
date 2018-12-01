import sqlite3,sys
from sqlite3 import Error

def create_connection(db_file):
    try:
        conn = sqlite3.connect(db_file)
        return conn
    except Error as e:
        print(e)
 
    return None
 
 
def select_all(conn):
    cur = conn.cursor()
    cur.execute("SELECT * FROM log")
    rows = cur.fetchall()
    for row in rows:
        print(row)

def main():

    database = "databases/" + sys.argv[1]
    print(database)
 
    conn = create_connection(database)
    with conn:
        select_all(conn)
 
if __name__ == '__main__':
    main()