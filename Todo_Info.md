# TO DO Info

## Arduino

### Setup function

For the setup function, just call an external function (called `startup` for example) that publishes to the `/startup` channel and awaits a response.

Set a boolean global variable indicating the state of the device (`initialized` -> `true`/`false`).

When modifying a channel component, call a `restart` function that unsubscribes from existent channels and then calls the `startup` function to reset those information and resubscribe to correct channels.

```c
void startup() {
  // Connect to MQTT
  // Send `/startup` message to MQTT Broker
  // Wait for MQTT Broker response
  // Set MQTT channels and other variables
  // Subscribe to MQTT channels
}

void setup() {
  Serial.begin(9600);

  WiFi.mode(WIFI_STA);
  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);

  Serial.println("ESP32 - Connecting to Wi-Fi");

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println();
  
  // Call the startup() function:
  startup();
}

void reset() {
  // Unsubscribe to MQTT channels
  
  // Call the startup() function:
  startup();
}
```

## Production Deployment

### Set the environment variables for the systemd service

Times change and so do best practices.

The current best way to do this is to `run systemctl edit myservice`, which will create an override file for you or let you edit an existing one.

In normal installations this will create a directory `/etc/systemd/system/myservice.service.d`, and inside that directory create a file whose name ends in `.conf` (typically, `override.conf`), and in this file you can add to or override any part of the unit shipped by the distribution.

For instance, in a file `/etc/systemd/system/myservice.service.d/myenv.conf`:

```toml
[Service]
Environment="SECRET=pGNqduRFkB4K9C2vijOmUDa2kPtUhArN"
Environment="ANOTHER_SECRET=JP8YLOc2bsNlrGuD6LVTq7L36obpjzxd"
```

Also note that if the directory exists and is empty, your service will be disabled! If you don't intend to put something in the directory, ensure that it does not exist.