-- Insert dummy data into subscriptions table
INSERT INTO subscriptions (tier, price_per_month) 
VALUES 
    ('free', 0.00),
    ('business', 299000.00);  -- Assuming the price is 299,000 IDR for the business tier

-- Insert dummy data into payments table
INSERT INTO payments (user_id, payment_gateway, amount, currency, transaction_date, status, url) 
VALUES 
    (1, 'xendit', 299000.00, 'IDR', NOW(), 'completed', 'https://payment-gateway.com/payment/1'),
    (2, 'paypal', 299000.00, 'IDR', NOW(), 'pending', 'https://payment-gateway.com/payment/2'),
    (3, 'stripe', 299000.00, 'IDR', NOW(), 'failed', 'https://payment-gateway.com/payment/3');

-- Insert dummy data into user_subscriptions table
INSERT INTO user_subscriptions (user_id, subscription_id, duration, end_date, payment_id)
VALUES 
    (1, 2, 12, NOW() + INTERVAL '12 months', 1),  -- Business tier for 12 months
    (2, 2, 6, NULL, 2),    
    (3, 1, 0, NULL, 3);     