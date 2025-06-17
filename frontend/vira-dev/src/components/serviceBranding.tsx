import { Shield, Code } from 'lucide-react';

interface ServiceBrandingProps {
  service: 'vira-id' | 'vira-dev' | 'transition';
  className?: string;
}

export function ServiceBranding({ service, className = '' }: ServiceBrandingProps) {
  if (service === 'transition') {
    return (
      <div className={`flex items-center space-x-2 ${className}`}>
        <div className="flex items-center space-x-1 opacity-50 transition-opacity duration-1000">
          <Shield className="w-6 h-6 text-blue-400" />
          <span className="text-xl font-bold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
            Vira.ID
          </span>
        </div>
        <div className="text-white/40 transition-opacity duration-1000">×</div>
        <div className="flex items-center space-x-1 transition-opacity duration-1000">
          <Code className="w-6 h-6 text-purple-400" />
          <span className="text-xl font-bold bg-gradient-to-r from-purple-400 to-pink-400 bg-clip-text text-transparent">
            Vira.Dev
          </span>
        </div>
      </div>
    );
  }

  const config = {
    'vira-id': {
      icon: Shield,
      name: 'Vira.ID',
      gradient: 'from-blue-400 to-purple-400',
      iconColor: 'text-blue-400',
    },
    'vira-dev': {
      icon: Code,
      name: 'Vira.Dev',
      gradient: 'from-purple-400 to-pink-400',
      iconColor: 'text-purple-400',
    },
  };

  const { icon: Icon, name, gradient, iconColor } = config[service];

  return (
    <div className={`flex items-center space-x-2 animate-slide-in ${className}`}>
      <Icon className={`w-8 h-8 ${iconColor}`} />
      <span className={`text-2xl font-bold bg-gradient-to-r ${gradient} bg-clip-text text-transparent`}>
        {name}
      </span>
    </div>
  );
}