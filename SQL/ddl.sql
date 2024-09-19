-- Create Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phonenumber VARCHAR(15) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    tier VARCHAR(20) CHECK (tier IN ('free', 'business')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Create Locations Table to store user's saved locations
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    location_name VARCHAR(100),
    latitude DECIMAL(9,6) NOT NULL,
    longitude DECIMAL(9,6) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Emissions Table to store data of business emissions for Business Tier users
CREATE TABLE business_facilities (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    company_type VARCHAR(50),
    total_emission DECIMAL(10,2) NOT NULL,  -- Total carbon emissions of the business
    location_id INT REFERENCES locations(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Emissions Table to store data of business emissions for free Tier users
CREATE TABLE user_location (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    location_id INT REFERENCES locations(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Air_Quality Data Table for storing air quality data from third-party API
CREATE TABLE air_quality (
    id SERIAL PRIMARY KEY,
    location_id INT REFERENCES locations(id) ON DELETE CASCADE,
    aqi INT NOT NULL,  -- Air Quality Index
    co DECIMAL(5,2),    -- Carbon Monoxide
    no DECIMAL(5,2), -- Nitrogen Monoxide
    no2 DECIMAL(5,2),   -- Nitrogen Dioxide
    o3 DECIMAL(5,2),    -- Ozone
    so2 DECIMAL(5,2),   -- Sulfur Dioxide
    pm25 DECIMAL(5,2),  -- Particulate matter <2.5 micrometers
    pm10 DECIMAL(5,2),  -- Particulate matter <10 micrometers
    nh3 DECIMAL(5,2), -- Ammonia
    fetch_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Payments Table for storing payment transactions for Business Tier users
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    payment_gateway VARCHAR(50) NOT NULL,  -- e.g., 'xendit'
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) CHECK (status IN ('pending', 'completed', 'failed')) NOT NULL
);
