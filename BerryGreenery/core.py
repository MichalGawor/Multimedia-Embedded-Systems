import datetime
import psycopg2
import RPi.GPIO as GPIO
import serial
import time


GPIO.setmode(GPIO.BOARD)
GPIO.setup(36, GPIO.OUT, initial=GPIO.LOW)
GPIO.setup(38, GPIO.OUT, initial=GPIO.LOW)
GPIO.setup(40, GPIO.OUT, initial=GPIO.LOW)

PUMP_PIN = 40
COOLER_PIN = 38
LIGHT_PIN = 36

DB_NAME = "test"
USER = "pi"
HOST = "localhost"
PASSWORD = "qwerty"


class BerryGreenery():
    def __init__(self, dbname, user, host, password):
        self.day_time = 3600
        self.night_time = 3600
        self.desired_moisture = 30
        self.desired_temperature = 25
        self.cycle_start_time = time.time()
        self.day = True
        self.db_connection, self.db_cursor = self.db_connect(dbname, user, host, password)
        self.db_get_and_set_expected()
        self.read_values()
        


    def db_connect(self, dbname, user, host, password):
        try:
            connection_str = "dbname='"+dbname+"' user='"+user+"' host='"+host+"' password='"+password+"'"
            db_connection = psycopg2.connect(connection_str)
            db_cursor = db_connection.cursor()
            return db_connection, db_cursor
        except Exception as error:
            print("Could not connect to database")
            print(error)

    def db_insert_measurments(self, moisture, water_level, temperature):
        sql_querry = """INSERT INTO real_table(moisture, "water", temperature)
                        VALUES(%s, %s, %s)"""   
        self.db_cursor.execute(sql_querry, (moisture, water_level, temperature))
        self.db_connection.commit()
        print("Inserted data")

    def db_get_and_set_expected(self):
        sql_querry = """SELECT * FROM expected_table ORDER by id DESC LIMIT 1;"""
        self.db_cursor.execute(sql_querry)
        records = self.db_cursor.fetchall()
        records = records[0]
        self.day_time = 3600 * records[1]
        self.night_time = 3600 * records[2]
        self.desired_moisture = records[3]
        self.desired_temperature = records[4]
        
        
    def read_values(self):
        uno_message = int(serialConnection.readline())
        while uno_message != -204:
            uno_message = int(serialConnection.readline())
        if uno_message == -204:
            self.moisture = int(serialConnection.readline())
            self.water_level = int(serialConnection.readline())
            self.temperature = int(serialConnection.readline())
            print("moisture: " , self.moisture, '%')
            print("temperature: ", self.temperature, '*C')
            print("water level: ", self.water_level, '%')
            self.db_insert_measurments(self.moisture, self.water_level, self.temperature)
    
    def controller(self):
        while True:
            self.read_values()
            if self.desired_moisture - self.moisture >= 10:
                self.pump_water()
            if self.temperature - self.desired_temperature >= 2:
                self.cool_air_on()
            if self.desired_temperature - self.temperature >= 2:
                self.cool_air_off()
            self.day_night_cycle()
    
    def pump_water(self):
        GPIO.output(PUMP_PIN, 1)
        time.sleep(2)
        GPIO.output(PUMP_PIN, 0)
        
    def cool_air_on(self):
        GPIO.output(COOLER_PIN, 1)
    
    def cool_air_off(self):
        GPIO.output(COOLER_PIN, 0)
        
    def day_night_cycle(self):
        if self.day and time.time() - self.cycle_start_time > self.day_time:
            self.day = False
            self.cycle_start_time = time.time()
            GPIO.output(LIGHT_PIN, 0)
        elif not self.day and time.time() - self.cycle_start_time > self.night_time:
            self.day = True
            self.cycle_start_time = time.time()
            GPIO.output(LIGHT_PIN, 1)
            

if __name__ == "__main__":
    # serial connection do arduino
    serialConnection = serial.Serial("/dev/ttyACM0", 9600)
    serialConnection.baudrate = 9600
    try:
        app_instance = BerryGreenery(DB_NAME, USER, HOST, PASSWORD)
        app_instance.controller()

    except KeyboardInterrupt:
        print("BerryGreenery turned off")
    finally:
        GPIO.cleanup()
