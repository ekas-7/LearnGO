-- Insert sample categories
INSERT INTO categories (id, name, description) VALUES
    ('550e8400-e29b-41d4-a716-446655440001', 'Electronics', 'Electronic devices and accessories'),
    ('550e8400-e29b-41d4-a716-446655440002', 'Clothing', 'Fashion and apparel'),
    ('550e8400-e29b-41d4-a716-446655440003', 'Books', 'Books and educational materials'),
    ('550e8400-e29b-41d4-a716-446655440004', 'Home & Garden', 'Home decor and gardening supplies')
ON CONFLICT (id) DO NOTHING;

-- Insert sample admin user (password: admin123)
-- Note: In production, ensure password is properly hashed
INSERT INTO users (id, email, password, first_name, last_name, role) VALUES
    ('550e8400-e29b-41d4-a716-446655440010', 
     'admin@example.com', 
     '$2a$10$YourHashedPasswordHere', 
     'Admin', 
     'User', 
     'admin')
ON CONFLICT (email) DO NOTHING;

-- Insert sample products
INSERT INTO products (id, name, description, price, stock, category_id, image_url) VALUES
    ('660e8400-e29b-41d4-a716-446655440001', 'Laptop', 'High-performance laptop', 999.99, 10, '550e8400-e29b-41d4-a716-446655440001', 'https://example.com/laptop.jpg'),
    ('660e8400-e29b-41d4-a716-446655440002', 'T-Shirt', 'Cotton t-shirt', 19.99, 50, '550e8400-e29b-41d4-a716-446655440002', 'https://example.com/tshirt.jpg'),
    ('660e8400-e29b-41d4-a716-446655440003', 'Programming Book', 'Learn Go programming', 39.99, 25, '550e8400-e29b-41d4-a716-446655440003', 'https://example.com/book.jpg')
ON CONFLICT (id) DO NOTHING;
