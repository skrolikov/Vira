import { useState } from 'react';
import { Eye, EyeOff } from 'lucide-react';

interface InputFieldProps {
  type: string;
  placeholder: string;
  value: string;
  onChange: (value: string) => void;
  icon?: React.ReactNode;
  required?: boolean;
  className?: string;
}

export function InputField({ 
  type, 
  placeholder, 
  value, 
  onChange, 
  icon, 
  required = false,
  className = '' 
}: InputFieldProps) {
  const [showPassword, setShowPassword] = useState(false);
  const [isFocused, setIsFocused] = useState(false);

  const inputType = type === 'password' && showPassword ? 'text' : type;

  return (
    <div className={`relative group ${className}`}>
      <div className={`
        relative flex items-center glass-input border border-white/10 rounded-2xl
        transition-all duration-300
        ${isFocused ? 'border-purple-400/50 shadow-lg shadow-purple-400/10' : 'hover:border-white/20'}
      `}>
        {icon && (
          <div className="absolute left-4 text-white/60 group-focus-within:text-purple-400 transition-colors duration-300">
            {icon}
          </div>
        )}
        
        <input
          type={inputType}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
          placeholder={placeholder}
          required={required}
          className={`
            w-full px-4 py-4 bg-transparent text-white placeholder-white/40 rounded-2xl outline-none
            ${icon ? 'pl-12' : ''}
            ${type === 'password' ? 'pr-12' : ''}
          `}
        />
        
        {type === 'password' && (
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute right-4 text-white/60 hover:text-purple-400 transition-colors duration-300"
          >
            {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
          </button>
        )}
      </div>
      
      {/* Animated underline */}
      <div className={`
        absolute bottom-0 left-0 h-0.5 bg-gradient-to-r from-purple-400 to-pink-400 rounded-full
        transition-all duration-300 
        ${isFocused ? 'w-full' : 'w-0'}
      `} />
    </div>
  );
}