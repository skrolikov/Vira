import { createBrowserRouter } from "react-router-dom";
import { Register } from "../pages/register";
import { Auth } from "../pages/auth";
import { Home } from "../pages/home";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "/auth",
    element: <Auth />,
  },
  {
    path: "/register",
    element: <Register />,
  },
]);
