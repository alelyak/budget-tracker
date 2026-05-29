import { Link, useNavigate, useLocation } from "react-router-dom";
import { Group, Button, Text } from "@mantine/core";

function Navbar() {
  const navigate = useNavigate();
  const location = useLocation();
  const user = JSON.parse(localStorage.getItem("user") || "{}");

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    navigate("/login");
  };

  const isActive = (path) => location.pathname === path;

  return (
    <Group
      justify="space-between"
      style={{ height: "100%", padding: "0 20px" }}
    >
      <Group>
        <Text fw={700} size="lg">
          💰 Финансы
        </Text>
        <Button
          variant={isActive("/dashboard") ? "light" : "subtle"}
          color={isActive("/dashboard") ? "blue" : "gray"}
          component={Link}
          to="/dashboard"
        >
          📊 Дашборд
        </Button>
        <Button
          variant={isActive("/transactions") ? "light" : "subtle"}
          color={isActive("/transactions") ? "blue" : "gray"}
          component={Link}
          to="/transactions"
        >
          💰 Транзакции
        </Button>
        <Button
          variant={isActive("/categories") ? "light" : "subtle"}
          color={isActive("/categories") ? "blue" : "gray"}
          component={Link}
          to="/categories"
        >
          📂 Категории
        </Button>
      </Group>
      <Group>
        <Text size="sm" c="dimmed">
          👤 {user.email}
        </Text>
        <Button variant="outline" onClick={handleLogout} color="red">
          🚪 Выйти
        </Button>
      </Group>
    </Group>
  );
}

export default Navbar;
