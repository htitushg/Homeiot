INSERT INTO locations (id, name, type, created_at, updated_at) VALUES
       (1, 'Living Room', 'Room', NOW(), NOW()),
       (2, 'Kitchen', 'Room', NOW(), NOW()),
       (3, 'Garage', 'Outdoor', NOW(), NOW());

INSERT INTO devices (id, name, type, location_id, created_at, updated_at) VALUES
      ('dev-001', 'Living Room Light', 'light', 1, NOW(), NOW()),
      ('dev-002', 'Kitchen Light', 'light', 2, NOW(), NOW()),
      ('dev-003', 'Garage Light', 'light', 3, NOW(), NOW());

INSERT INTO modules (id, name, value, device_id, created_at, updated_at) VALUES
     (1, 'temperatureSensor', '22.5', 'dev-001', NOW(), NOW()),
     (2, 'luminositySensor', '200', 'dev-001', NOW(), NOW()),
     (3, 'presenceDetector', 'True', 'dev-002', NOW(), NOW()),
     (4, 'consumptionSensor', '150', 'dev-002', NOW(), NOW()),
     (5, 'lightSensor', 'False', 'dev-003', NOW(), NOW());

INSERT INTO data (id, module_id, module_name, module_value, device_id, created_at, updated_at) VALUES
   (1, 1, 'temperatureSensor', '22.5', 'dev-001', NOW(), NOW()),
   (2, 2, 'luminositySensor', '200', 'dev-001', NOW(), NOW()),
   (3, 3, 'presenceDetector', 'True', 'dev-002', NOW(), NOW()),
   (4, 4, 'consumptionSensor', '150', 'dev-002', NOW(), NOW()),
   (5, 5, 'lightSensor', 'False', 'dev-003', NOW(), NOW());