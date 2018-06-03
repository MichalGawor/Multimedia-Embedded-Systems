#include <DHT.h>

#include <Adafruit_Sensor.h>

#define DHT_PIN 3
#define MOISTURE_SENSOR_PIN A0
#define WATER_LEVEL_SENSOR_PIN A2
#define DHT_TYPE DHT11


DHT dht(DHT_PIN, DHT_TYPE);

void setup(){
  // TODO 
  // WATER_LEVEL_SENSOR_PIN = 
  // MOISTURE_SENSOR_PIN = 
  // initialize serial communication at 9600 bits per second:
  Serial.begin(9600);
  dht.begin();
}


void readSensors(int n, int _delay){
  int averageMoisture = 0;
  int averageWaterLevel = 0;
  int averageTemperature = 0;
  for(int i = 0; i<n; i++){
    delay(_delay);
    averageMoisture += analogRead(MOISTURE_SENSOR_PIN);
    averageWaterLevel += analogRead(WATER_LEVEL_SENSOR_PIN);
    averageTemperature += dht.readTemperature();
  }
  averageMoisture /= n;
  averageWaterLevel /= n;
  averageTemperature /= n;
  averageMoisture = map(averageMoisture, 1024, 0, 0, 100);
  averageWaterLevel = map(averageWaterLevel, 0, 750, 0, 100);
  Serial.print(-204);
  Serial.print('\n');
  Serial.print(averageMoisture);
  Serial.print('\n');
  Serial.print(averageWaterLevel);
  Serial.print('\n');
  Serial.print(averageTemperature);
  Serial.print('\n');
}

void loop() {
  // Read sensor values
  // int waterLevelValue = analogRead(WATER_LEVEL_SENSOR_PIN);

  // Compute numerical value
  // readTest();
  readSensors(2, 500);
  delay(1000);
  
}



