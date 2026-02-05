# Flight control

https://claude.ai/chat/f9311a7c-91c9-4b05-baf9-b8da483ed746


## Problem

The mobile could replace the raspberry pi and talk to the Pixhawk gear ?

Basically we can install our App on Mobiles and it is the sensor package with:

- coms
- camera
- gps
- IMU

Mobiles are easiyl the most secure devices around.

And out System can use NATS for sync to cloud and Operators own Mobile / Desktop / Server

Goup-util just needs to support the tricks needed to keep the App as Foregroudn as possible.



## Solutions

## Foregound

The App Store Approval Game
Apple reviews apps for legitimate use of background modes. You can't just tick all the boxes - they'll reject you.

What gets approved:

Background ModeApple ExpectsYour JustificationLocationNavigation, fitness, tracking✅ "Drone flight logging & geofencing"External AccessoryMFi hardware⚠️ Need MFi or BluetoothAudioMusic, podcasts, voice✅ "Audible flight alerts & warnings"BluetoothConnected peripherals✅ "Drone telemetry receiver"Background fetchContent updates✅ "Flight plan sync"


The Winning Narrative for Apple Review
Your app isn't a "hack to stay alive" - it's a legitimate drone companion app:

App Description:
"DroneLink Pro - Professional UAV Ground Station

Real-time telemetry monitoring and flight management for 
MAVLink-compatible drones. Features:

- Live GPS tracking and geofence alerts
- Audible warnings for battery, altitude, signal loss  
- Bluetooth telemetry integration
- Flight logging and blackbox recording
- Store-and-forward for areas with poor connectivity"
Every background mode has a user-facing reason.
Technical Architecture That Passes Review
swift// Info.plist - all justified
<key>UIBackgroundModes</key>
<array>
    <string>location</string>           // Flight tracking, geofence
    <string>bluetooth-central</string>  // Telemetry from drone
    <string>audio</string>              // Voice alerts
    <string>fetch</string>              // Sync flight plans
</array>

<key>NSLocationAlwaysAndWhenInUseUsageDescription</key>
<string>Track your drone's position and enforce geofence boundaries</string>

<key>NSBluetoothAlwaysUsageDescription</key>  
<string>Receive real-time telemetry from your drone</string>
```

## The Bluetooth Bridge Solution

Since USB is locked down, use **Bluetooth Serial** to talk to Pixhawk:
```
┌──────────────┐     Bluetooth     ┌─────────────────┐
│   iPhone     │◄──────────────────►│  HC-05/HM-10    │
│   Your App   │    SPP or BLE      │  Serial Module  │
└──────────────┘                    └────────┬────────┘
                                             │ UART
                                    ┌────────▼────────┐
                                    │    Pixhawk      │
                                    │   TELEM2 Port   │
                                    └─────────────────┘
Hardware cost: ~$5 for HC-05 module
Apps Already Doing This (Proof It Works)
These are on the App Store RIGHT NOW:
AppWhat It DoesBackground ModesQGroundControlMAVLink ground stationLocation, BluetoothDJI FlyDrone controlLocation, Bluetooth, AudioLitchiAutonomous flightsLocation, BluetoothDroneLinkMission planningLocation, Bluetooth
Apple approves drone apps. You just need the right framing.
Your App Store Strategy
Phase 1: MVP that gets approved

Bluetooth telemetry viewer
GPS flight logging
Basic alerts (audio)
This WILL pass review

Phase 2: Add features

NATS connectivity for remote monitoring
Store-and-forward telemetry
Advanced sensor fusion
Apple rarely re-reviews existing approved background modes

Code Structure (Swift + Your Go Backend)
swift// DroneLink iOS App

class FlightManager: NSObject, CLLocationManagerDelegate, CBCentralManagerDelegate {
    
    let locationManager = CLLocationManager()
    let bluetoothManager: CBCentralManager
    let audioEngine = AVAudioEngine()  // For alerts
    var natsConnection: NATSClient?    // Your Go server
    
    func startFlight() {
        // Location - keeps app alive
        locationManager.startUpdatingLocation()
        locationManager.allowsBackgroundLocationUpdates = true
        
        // Bluetooth - telemetry from drone
        bluetoothManager.scanForPeripherals(withServices: [mavlinkServiceUUID])
        
        // Audio - alert system (also keeps alive)
        startAlertAudioSession()
        
        // NATS - your backend
        natsConnection?.connect(to: "nats://your-server.com")
    }
    
    func locationManager(_ manager: CLLocationManager, 
                         didUpdateLocations locations: [CLLocation]) {
        // This fires continuously in background
        // Perfect place to:
        // - Check geofence
        // - Log flight path  
        // - Send position to NATS
        // - Process incoming MAVLink
    }
}
Your Backend (Go + NATS)
Fits perfectly with your existing stack:
go// Ground station server
type GroundStation struct {
    js   nats.JetStreamContext
    pb   *pocketbase.PocketBase  // Flight logs, auth
}

func (gs *GroundStation) HandleTelemetry(msg *nats.Msg) {
    // Phone sends telemetry here
    // Store in PocketBase
    // Push to web dashboard via Datastar SSE
}
Summary: Yes, You Can Ship This

Use Bluetooth instead of USB for Pixhawk comms
Frame it as drone ground station - legitimate use case
Every background mode has a user-facing feature
QGroundControl proves Apple approves this category
Your Go/NATS backend handles the heavy lifting

Want me to sketch out the full iOS app architecture, or the Bluetooth-MAVLink bridge hardware setup?



## Bluetooth-MAVLink Bridge Hardware Setup

You have **three good options** depending on your needs:

---

## Option 1: ESP32 DroneBridge (Best for iOS)

**Why:** WiFi works perfectly with iOS, no MFi needed. BLE also supported.

### Parts List (~$15-25)

| Part | Cost | Notes |
|------|------|-------|
| ESP32 Dev Board (with antenna) | $8-15 | ESP32-WROOM or WROVER |
| DF13 6-pin cable | $2 | For Pixhawk TELEM port |
| Optional: U.FL external antenna | $5 | For extended range |

### Wiring

```
┌─────────────────────────────────────────────────────────┐
│                    PIXHAWK                              │
│                                                         │
│   TELEM1 or TELEM2 Port (DF13 6-pin)                   │
│   ┌───┬───┬───┬───┬───┬───┐                            │
│   │5V │TX │RX │CTS│RTS│GND│                            │
│   └─┬─┴─┬─┴─┬─┴───┴───┴─┬─┘                            │
│     │   │   │           │                               │
└─────┼───┼───┼───────────┼───────────────────────────────┘
      │   │   │           │
      │   │   │           │
┌─────▼───▼───▼───────────▼───────────────────────────────┐
│     VIN  RX  TX        GND                              │
│                                                         │
│              ESP32 Dev Board                            │
│         (Running DroneBridge firmware)                  │
│                                                         │
│    WiFi AP: "DroneBridge ESP32"                        │
│    Password: "dronebridge"                              │
│    TCP Port: 5760 / UDP Port: 14550                    │
└─────────────────────────────────────────────────────────┘
```

**Important:** TX→RX, RX→TX (crossed!)

### Flashing DroneBridge

```bash
# Download from https://github.com/DroneBridge/ESP32
# Use ESP Flash Tool or esptool.py

esptool.py --chip esp32 --port /dev/ttyUSB0 \
  write_flash 0x0 dronebridge_esp32.bin
```

### iOS Connection

```swift
// Your app connects via WiFi - no Bluetooth needed!
let socket = GCDAsyncUdpSocket(delegate: self, delegateQueue: .main)
try socket.bind(toPort: 14550)
try socket.beginReceiving()

// Send heartbeat to register with ESP32
let heartbeat = MAVLinkHeartbeat()
socket.send(heartbeat.data, toHost: "192.168.2.1", port: 14550, withTimeout: -1, tag: 0)
```

---

## Option 2: HC-05/HC-06 Classic Bluetooth (Android Only)

**Why:** Cheap, simple, works with SPP (Serial Port Profile)

**⚠️ Won't work with iOS** - iOS doesn't support Bluetooth SPP

### Parts List (~$5-10)

| Part | Cost |
|------|------|
| HC-05 or HC-06 module | $4-6 |
| DF13 to Dupont cables | $2 |

### Wiring

```
┌─────────────────────────────────────────────────────────┐
│                    PIXHAWK                              │
│                                                         │
│   TELEM2 Port (DF13 6-pin)                             │
│   ┌───┬───┬───┬───┬───┬───┐                            │
│   │5V │TX │RX │CTS│RTS│GND│                            │
│   └─┬─┴─┬─┴─┬─┴───┴───┴─┬─┘                            │
│     │   │   │           │                               │
└─────┼───┼───┼───────────┼───────────────────────────────┘
      │   │   │           │
      │   │   │           │
┌─────▼───▼───▼───────────▼───────────────────────────────┐
│    VCC  RX  TX         GND                              │
│     ●    ●   ●    ●     ●    ●                         │
│    VCC  GND TXD  RXD  STATE KEY                        │
│                                                         │
│              HC-05 Bluetooth Module                     │
│                                                         │
│    Default name: HC-05 or HC-06                        │
│    Pairing PIN: 1234 or 0000                           │
│    Default baud: 9600 (change to 57600)                │
└─────────────────────────────────────────────────────────┘
```

### Configure HC-05 Baud Rate (One-time)

```
# Connect HC-05 to FTDI USB-Serial adapter
# Hold KEY button while powering on (enters AT mode)
# Open serial terminal at 38400 baud

AT                    # Should respond "OK"
AT+NAME=DroneLink     # Set device name
AT+UART=57600,0,0     # Set baud to match Pixhawk
AT+PSWD=1234          # Set PIN
```

---

## Option 3: HM-10 BLE (iOS Compatible via Custom App)

**Why:** BLE 4.0 works with iOS Core Bluetooth

**Caveat:** Need to implement GATT service in your app, QGroundControl won't auto-connect

### Parts List (~$6-12)

| Part | Cost |
|------|------|
| HM-10 BLE module | $5-8 |
| DF13 cables | $2 |

### Wiring (Same as HC-05)

```
Pixhawk TELEM2          HM-10
─────────────          ───────
    5V  ──────────────► VCC
    TX  ──────────────► RXD
    RX  ◄──────────────  TXD
   GND  ──────────────► GND
```

### iOS Code for HM-10

```swift
import CoreBluetooth

// HM-10 Service & Characteristic UUIDs
let HM10_SERVICE = CBUUID(string: "FFE0")
let HM10_CHAR = CBUUID(string: "FFE1")

class DroneBLE: NSObject, CBCentralManagerDelegate, CBPeripheralDelegate {
    var central: CBCentralManager!
    var peripheral: CBPeripheral?
    var characteristic: CBCharacteristic?
    
    func centralManagerDidUpdateState(_ central: CBCentralManager) {
        if central.state == .poweredOn {
            central.scanForPeripherals(withServices: [HM10_SERVICE])
        }
    }
    
    func peripheral(_ peripheral: CBPeripheral, 
                    didUpdateValueFor characteristic: CBCharacteristic, error: Error?) {
        guard let data = characteristic.value else { return }
        // Parse MAVLink from data
        processMAVLink(data)
    }
    
    func sendMAVLink(_ message: Data) {
        guard let char = characteristic else { return }
        peripheral?.writeValue(message, for: char, type: .withoutResponse)
    }
}
```

---

## My Recommendation for Your Use Case

Given you want **iOS App Store distribution**:

### Use ESP32 DroneBridge with WiFi

```
┌────────────────────────────────────────────────────────────┐
│                     ON THE DRONE                           │
│                                                            │
│   ┌──────────────┐      UART      ┌──────────────────┐    │
│   │   Pixhawk    │◄──────────────►│  ESP32 Module    │    │
│   │              │  (57600 baud)  │  (DroneBridge)   │    │
│   │  TELEM1/2    │                │                  │    │
│   └──────────────┘                └────────┬─────────┘    │
│                                            │ WiFi AP      │
└────────────────────────────────────────────┼──────────────┘
                                             │
                                             │ 192.168.2.1
                                             │ UDP 14550
                                             ▼
┌────────────────────────────────────────────────────────────┐
│                     YOUR iPHONE APP                        │
│                                                            │
│   ┌──────────────────────────────────────────────────┐    │
│   │  Connect to WiFi: "DroneBridge ESP32"            │    │
│   │  UDP Socket → 192.168.2.1:14550                  │    │
│   │                                                  │    │
│   │  Background modes:                               │    │
│   │  ✓ Location (keeps alive)                        │    │
│   │  ✓ Audio (voice alerts)                          │    │
│   │  ✓ Background fetch                              │    │
│   │                                                  │    │
│   │  Also connects to LTE for:                       │    │
│   │  → NATS telemetry to your server                │    │
│   │  → Store-and-forward when WiFi drops            │    │
│   └──────────────────────────────────────────────────┘    │
│                                                            │
└────────────────────────────────────────────────────────────┘
```

### Why This Works for App Store

1. **WiFi connection is standard** - no special entitlements needed
2. **Background location** - legitimate for flight tracking
3. **Background audio** - legitimate for flight alerts
4. **No MFi required** - WiFi, not USB/Bluetooth accessory
5. **QGroundControl does exactly this** - proven App Store approved

### Parts to Order

| Item | Link/Search | Cost |
|------|-------------|------|
| ESP32-WROOM-32U (external antenna) | AliExpress/Amazon | $8 |
| 5dBi 2.4GHz antenna + U.FL pigtail | AliExpress | $3 |
| DF13 6-pin to Dupont cable | AliExpress | $2 |
| **Total** | | **~$13** |

Want me to sketch out the Swift app architecture that bridges ESP32 WiFi ↔ your NATS backend?


### WILD STUFF going on


https://wildhornets.com/en/

This is the heavy hitter of the volunteer-to-industry pipeline. Their site is where they showcase the Sting interceptors and the Queen Hornet (the "Mavic replacement" bomber).
What to look for: Check their "Products" or "News" section for the 2026 updates on the Hornet Vision digital link system—it’s the tech they built specifically to compete with DJI’s O3 Air Unit but without the Chinese telemetry.
Twitter/X: @wilendhornets (Highly active with footage of their drones taking down Shaheds).

---

https://firstcontact.biz/en/

These guys have been at it since 2014, long before the current craze. They are a formal MilTech firm rather than just a volunteer group.
The "Yautja" (Predator): On their projects page, look for the Yautja and Vidsich models. These are the hardened, non-Chinese quads designed for recon and strike in high-EW environments.
Deep Tech: They also have some wild projects on their site like the VAL-1, an underwater kamikaze drone with an 800km range.

---

https://militarnyi.com/en/

If you want to track the technical teardowns of these drones (comparing them to the Mav 3 you flew), this is the "Wall Street Journal" of Ukrainian defense tech. They do deep dives into the flight controllers and radio links of the First Contact and Wild Hornet fleets.

A Pro Tip for the Recon Pilot: > If you're looking for the software side of how they replaced DJI Fly, keep an eye on "Mission Control" and "Delta"—the Ukrainian battlefield management systems that these drones now plug into directly.

---

https://github.com/o-gs/dji-firmware-tools

PYTHON

---

https://github.com/ANG13T/DroneXtract

GOLANG

DroneXtract is a comprehensive digital forensics suite for DJI drones made with Golang. It can be used to analyze drone sensor values and telemetry data, visualize drone flight maps, audit for criminal activity, and extract pertinent data within multiple file formats.