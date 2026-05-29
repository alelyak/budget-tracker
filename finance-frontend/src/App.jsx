import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AppShell } from "@mantine/core";
import Login from "./pages/Login";
import Dashboard from "./pages/Dashboard";
import Transactions from "./pages/Transactions";
import Categories from "./pages/Categories";
import Navbar from "./components/Navbar";
import ProtectedRoute from "./components/ProtectedRoute";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <AppShell header={{ height: 60 }}>
                <AppShell.Header>
                  <Navbar />
                </AppShell.Header>
                <AppShell.Main p="md" pt="70px">
                  <Dashboard />
                </AppShell.Main>
              </AppShell>
            </ProtectedRoute>
          }
        />
        <Route
          path="/transactions"
          element={
            <ProtectedRoute>
              <AppShell header={{ height: 60 }}>
                <AppShell.Header>
                  <Navbar />
                </AppShell.Header>
                <AppShell.Main p="md" pt="70px">
                  <Transactions />
                </AppShell.Main>
              </AppShell>
            </ProtectedRoute>
          }
        />
        <Route
          path="/categories"
          element={
            <ProtectedRoute>
              <AppShell header={{ height: 60 }}>
                <AppShell.Header>
                  <Navbar />
                </AppShell.Header>
                <AppShell.Main p="md" pt="70px">
                  <Categories />
                </AppShell.Main>
              </AppShell>
            </ProtectedRoute>
          }
        />
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
