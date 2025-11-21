-- Mock hazard zones around JP Nagar, Bengaluru (12.914599, 77.595036)
-- For heat map visualization within 5km radius
-- Run this in TiDB Cloud Chat2Query

USE test;

-- Clear existing data (optional)
DELETE FROM hazard_zones;

-- Insert mock hazard zones with varying crime rates and safety levels
-- Coordinates are within ~5km of JP Nagar center

INSERT INTO `hazard_zones` (`name`, `latitude`, `longitude`, `radius_km`, `severity`, `safe_rating`, `description`) VALUES
-- High crime areas (safe_rating 1.0-2.5)
('JP Nagar Market Backlanes', 12.910000, 77.590000, 0.3, 'high', 2.1, 'Poorly lit market backlanes. Avoid after 9pm. Multiple petty theft incidents reported.'),
('Bannerghatta Road Underpass', 12.900000, 77.600000, 0.4, 'high', 1.8, 'Dark underpass with limited visibility. High risk after sunset. Use main road instead.'),
('BTM Layout Industrial Area', 12.920000, 77.610000, 0.5, 'high', 2.3, 'Deserted industrial area at night. Low foot traffic increases risk. Stay on well-lit streets.'),

-- Medium risk areas (safe_rating 2.6-3.5)
('Jayanagar 4th Block Park Perimeter', 12.930000, 77.580000, 0.25, 'medium', 3.2, 'Park perimeter has limited lighting. Stick to main roads after 10pm.'),
('Banashankari Bus Stand Area', 12.905000, 77.575000, 0.35, 'medium', 3.0, 'Busy during day but isolated at night. Be cautious after 11pm.'),
('HSR Layout Outer Ring Road', 12.940000, 77.640000, 0.4, 'medium', 3.4, 'Fast-moving traffic area. Pedestrian accidents reported. Use designated crossings.'),

-- Lower risk but still caution areas (safe_rating 3.6-4.0)
('JP Nagar Phase 7 Construction Zone', 12.915000, 77.605000, 0.2, 'medium', 3.8, 'Ongoing construction. Uneven surfaces and limited lighting. Exercise caution.'),
('Bilekahalli Main Road', 12.925000, 77.620000, 0.3, 'low', 3.9, 'Generally safe but occasional incidents. Stay alert during late hours.'),

-- Safe areas (for contrast - safe_rating 4.1-5.0) - these won't trigger nudges
('JP Nagar Metro Station Area', 12.914599, 77.595036, 0.15, 'low', 4.8, 'Well-lit metro station area. High foot traffic. Generally safe.'),
('Forum Mall Vicinity', 12.920000, 77.590000, 0.2, 'low', 4.9, 'Commercial area with security presence. Very safe zone.');

-- Also update weather alerts to be near your location
DELETE FROM weather_alerts;
INSERT INTO `weather_alerts` (`latitude`, `longitude`, `radius_km`, `type`, `description`, `severity`) VALUES
(12.914599, 77.595036, 5.00, 'rain', 'Heavy monsoon showers expected in next 2 hours. Carry umbrella and avoid waterlogged areas.', 'medium'),
(12.914599, 77.595036, 3.00, 'fog', 'Early morning fog expected tomorrow 6-8 AM. Reduced visibility for driving.', 'low');

