-- Crear la base de datos si no existe
CREATE DATABASE IF NOT EXISTS burrowdb;
USE burrowdb;

-- Crear la tabla users si no existe
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    age INT,
    city VARCHAR(100),
    country VARCHAR(100),
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

-- Insertar datos de ejemplo solo si la tabla está vacía
INSERT IGNORE INTO users (name, email, age, city, country, phone, is_active) VALUES
('Federico Pereira', 'federico.pereira@email.com', 28, 'Buenos Aires', 'Argentina', '+54-11-1234-5678', TRUE),
('María García', 'maria.garcia@email.com', 34, 'Madrid', 'España', '+34-91-123-4567', TRUE),
('Carlos López', 'carlos.lopez@email.com', 45, 'México DF', 'México', '+52-55-1234-5678', TRUE),
('Ana Martínez', 'ana.martinez@email.com', 29, 'Bogotá', 'Colombia', '+57-1-234-5678', TRUE),
('Luis Rodríguez', 'luis.rodriguez@email.com', 38, 'Lima', 'Perú', '+51-1-234-5678', TRUE),
('Elena Fernández', 'elena.fernandez@email.com', 31, 'Barcelona', 'España', '+34-93-123-4567', TRUE),
('Miguel Santos', 'miguel.santos@email.com', 42, 'São Paulo', 'Brasil', '+55-11-1234-5678', TRUE),
('Carmen Díaz', 'carmen.diaz@email.com', 26, 'Santiago', 'Chile', '+56-2-1234-5678', TRUE),
('Roberto Silva', 'roberto.silva@email.com', 33, 'Montevideo', 'Uruguay', '+598-2-123-4567', TRUE),
('Sofía Morales', 'sofia.morales@email.com', 27, 'Quito', 'Ecuador', '+593-2-123-4567', TRUE),
('Fernando Castro', 'fernando.castro@email.com', 39, 'La Paz', 'Bolivia', '+591-2-123-4567', TRUE),
('Isabel Herrera', 'isabel.herrera@email.com', 35, 'Caracas', 'Venezuela', '+58-212-123-4567', TRUE),
('Diego Vargas', 'diego.vargas@email.com', 41, 'Asunción', 'Paraguay', '+595-21-123-456', TRUE),
('Lucía Jiménez', 'lucia.jimenez@email.com', 30, 'San José', 'Costa Rica', '+506-2123-4567', TRUE),
('Andrés Ruiz', 'andres.ruiz@email.com', 37, 'Panamá', 'Panamá', '+507-123-4567', TRUE),
('Patricia Romero', 'patricia.romero@email.com', 32, 'Tegucigalpa', 'Honduras', '+504-123-4567', TRUE),
('José Guerrero', 'jose.guerrero@email.com', 44, 'San Salvador', 'El Salvador', '+503-123-4567', TRUE),
('Laura Mendoza', 'laura.mendoza@email.com', 25, 'Guatemala', 'Guatemala', '+502-123-4567', TRUE),
('Manuel Ortega', 'manuel.ortega@email.com', 36, 'Managua', 'Nicaragua', '+505-123-4567', TRUE),
('Victoria Ramos', 'victoria.ramos@email.com', 43, 'Córdoba', 'Argentina', '+54-351-123-4567', TRUE);

-- Crear usuario con permisos completos
-- Primero eliminamos el usuario si existe para evitar conflictos
DROP USER IF EXISTS 'burrowuser'@'%';
DROP USER IF EXISTS 'burrowuser'@'localhost';

-- Crear el usuario desde cualquier host
CREATE USER 'burrowuser'@'%' IDENTIFIED BY 'burrowpass123';
CREATE USER 'burrowuser'@'localhost' IDENTIFIED BY 'burrowpass123';

-- Otorgar TODOS los privilegios en la base de datos burrowdb
GRANT ALL PRIVILEGES ON burrowdb.* TO 'burrowuser'@'%';
GRANT ALL PRIVILEGES ON burrowdb.* TO 'burrowuser'@'localhost';

-- Permisos adicionales para operaciones específicas
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, INDEX, ALTER ON burrowdb.* TO 'burrowuser'@'%';
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, INDEX, ALTER ON burrowdb.* TO 'burrowuser'@'localhost';

-- Aplicar los cambios
FLUSH PRIVILEGES;

-- Verificaciones finales
SELECT 'Tabla users creada exitosamente' as message;
SELECT COUNT(*) as total_users FROM users;

-- Mostrar los usuarios creados
SELECT 'Usuarios de MariaDB creados:' as message;
SELECT User, Host FROM mysql.user WHERE User = 'burrowuser';

-- Mostrar permisos del usuario
SELECT 'Permisos del usuario burrowuser:' as message;
SHOW GRANTS FOR 'burrowuser'@'%';

-- Verificar conexión del usuario (esto ayuda a validar que puede conectarse)
SELECT 'Usuario burrowuser puede consultar la tabla users:' as message;
SELECT id, name FROM users LIMIT 3; 