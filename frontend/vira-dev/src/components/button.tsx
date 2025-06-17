import { Loader2 } from 'lucide-react';
import clsx from 'clsx';

interface ButtonProps {
  children: React.ReactNode;
  onClick?: () => void;
  type?: 'button' | 'submit';
  variant?: 'primary' | 'secondary' | 'outline';
  size?: 'sm' | 'md' | 'lg';
  loading?: boolean;
  disabled?: boolean;
  className?: string;
  icon?: React.ReactNode;
}

export function Button({
  children,
  onClick,
  type = 'button',
  variant = 'primary',
  size = 'md',
  loading = false,
  disabled = false,
  className = '',
  icon,
}: ButtonProps) {
  const baseClasses =
    'relative inline-flex items-center justify-center overflow-hidden rounded-xl font-medium focus:outline-none focus:ring-2 focus:ring-purple-400/50 disabled:opacity-50 disabled:cursor-not-allowed transition-all group';

  const variants = {
    primary:
      'bg-purple-600 text-white hover:bg-purple-700 active:bg-purple-800',
    secondary:
      'bg-white/10 text-white border border-white/20 hover:bg-white/20',
    outline:
      'border-2 border-purple-500 text-purple-500 hover:bg-purple-500 hover:text-white',
  };

  const sizes = {
    sm: 'px-4 py-2 text-sm',
    md: 'px-6 py-3 text-base',
    lg: 'px-8 py-4 text-lg',
  };

  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled || loading}
      className={clsx(
        baseClasses,
        variants[variant],
        sizes[size],
        className
      )}
    >
      {/* Анимированный фон */}
      <span className="absolute inset-0 bg-white/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none" />

      {/* Контент */}
      <span className="relative flex items-center space-x-2">
        {loading ? (
          <Loader2 className="w-5 h-5 animate-spin" />
        ) : icon ? (
          <span>{icon}</span>
        ) : null}
        <span>{children}</span>
      </span>
    </button>
  );
}
