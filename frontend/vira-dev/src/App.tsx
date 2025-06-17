import { AnimatedBackground } from './components/animatedBackground'
import { RouterProvider } from 'react-router-dom'
import { router } from './router'
import { ThemeProvider } from './context/themeContext';

function App() {
  return (
    <ThemeProvider>
    <div className="min-h-screen transition-all duration-1000">
      <AnimatedBackground />

      <RouterProvider router={router} />
    </div>
    </ThemeProvider>
  );
}

export default App;