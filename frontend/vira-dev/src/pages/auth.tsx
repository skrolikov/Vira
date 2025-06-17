import { useState } from 'react';
import { Mail, Lock, ArrowRight, Shield, Code, CheckCircle} from 'lucide-react';

import { useNavigate } from 'react-router-dom';
import { GlassCard } from '../components/glassCard';
import { InputField } from '../components/inputField';
import { Button } from '../components/button';
import { ServiceBranding } from '../components/serviceBranding';

type AuthStep = 'vira-id' | 'transition' | 'vira-dev';

export function Auth() {
  const [step, setStep] = useState<AuthStep>('vira-id');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [isTransitioning, setIsTransitioning] = useState(false);
  const [rememberMe, setRememberMe] = useState(false);
  const navigate = useNavigate();

  const handleViraIdAuth = async () => {
    setLoading(true);
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1500));
    setLoading(false);
    
    // Start transition animation
    setIsTransitioning(true);
    setStep('transition');
    
    // After transition animation, move to Vira.Dev step
    setTimeout(() => {
      setStep('vira-dev');
      setIsTransitioning(false);
    }, 2000);
  };

  const handleViraDevAuth = async () => {
    setLoading(true);
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1500));
    setLoading(false);
    navigate('/dashboard');
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4 relative overflow-hidden">
      {/* Enhanced background effects */}
      <div className="absolute inset-0 bg-gradient-to-br from-purple-900/20 via-transparent to-pink-900/20" />
      
      <div className="w-full max-w-md relative z-10">
        {/* Service Branding with enhanced transitions */}
        <div className="text-center mb-8">
          <div className={`transition-all duration-1000 ease-out ${
            step === 'transition' 
              ? 'opacity-100 transform scale-110' 
              : isTransitioning 
                ? 'opacity-0 transform scale-95' 
                : 'opacity-100 transform scale-100'
          }`}>
            <ServiceBranding
              service={step === 'vira-dev' ? 'vira-dev' : step === 'transition' ? 'transition' : 'vira-id'} 
              className="justify-center mb-6"
            />
            
            {step !== 'transition' && (
              <div className={`transition-all duration-800 ${isTransitioning ? 'opacity-0 translate-y-4' : 'opacity-100 translate-y-0'}`}>
                <h1 className="text-3xl font-bold text-white mb-3">
                  {step === 'vira-id' ? 'Добро пожаловать' : 'Почти готово!'}
                </h1>
                <p className="text-white/70 text-lg leading-relaxed">
                  {step === 'vira-id' 
                    ? 'Войдите в вашу межсервисную учетную запись для доступа ко всем сервисам экосистемы Vira' 
                    : 'Ваш аккаунт Vira.ID успешно подтвержден. Сейчас мы создадим для вас профиль в образовательной платформе Vira.Dev'
                  }
                </p>
              </div>
            )}
          </div>
        </div>

        {/* Transition Animation */}
        {step === 'transition' && (
          <div className="text-center py-12">
            <div className="relative">
              <div className="w-20 h-20 mx-auto mb-6 rounded-full bg-gradient-to-br from-purple-500/20 to-pink-500/20 flex items-center justify-center animate-pulse">
                <CheckCircle className="w-10 h-10 text-green-400 animate-bounce" />
              </div>
              <h2 className="text-2xl font-bold text-white mb-4">Аутентификация успешна!</h2>
              <p className="text-white/70 mb-6">Подключаемся к Vira.Dev...</p>
              
              {/* Loading animation */}
              <div className="flex justify-center space-x-1">
                {[0, 1, 2].map((i) => (
                  <div
                    key={i}
                    className="w-2 h-2 bg-purple-400 rounded-full animate-bounce"
                    style={{ animationDelay: `${i * 0.2}s` }}
                  />
                ))}
              </div>
            </div>
          </div>
        )}

        {/* Auth Form */}
        {step !== 'transition' && (
          <GlassCard className={`p-8 transition-all duration-800 ${
            isTransitioning ? 'opacity-0 transform translate-x-8 scale-95' : 'opacity-100 transform translate-x-0 scale-100'
          }`}>
            <form onSubmit={(e) => {
              e.preventDefault();
              step === 'vira-id' ? handleViraIdAuth() : handleViraDevAuth();
            }} className="space-y-6">
              
              {step === 'vira-id' ? (
                <>
                  <div className="space-y-5">
                    <InputField
                      type="email"
                      placeholder="Введите ваш email"
                      value={email}
                      onChange={setEmail}
                      icon={<Mail size={20} />}
                      required
                    />
                    
                    <InputField
                      type="password"
                      placeholder="Введите пароль"
                      value={password}
                      onChange={setPassword}
                      icon={<Lock size={20} />}
                      required
                    />
                  </div>

                  {/* Remember me checkbox */}
                  <div className="flex items-center justify-between">
                    <label className="flex items-center space-x-3 cursor-pointer group">
                      <div className="relative">
                        <input
                          type="checkbox"
                          checked={rememberMe}
                          onChange={(e) => setRememberMe(e.target.checked)}
                          className="sr-only"
                        />
                        <div className={`w-5 h-5 rounded-lg border-2 transition-all duration-300 ${
                          rememberMe 
                            ? 'bg-purple-500 border-purple-500' 
                            : 'border-white/30 group-hover:border-purple-400'
                        }`}>
                          {rememberMe && (
                            <CheckCircle className="w-3 h-3 text-white absolute top-0.5 left-0.5" />
                          )}
                        </div>
                      </div>
                      <span className="text-white/70 group-hover:text-white transition-colors duration-300">
                        Запомнить меня
                      </span>
                    </label>
                    
                    <button 
                      type="button"
                      className="text-purple-400 hover:text-purple-300 transition-colors duration-300 text-sm"
                    >
                      Забыли пароль?
                    </button>
                  </div>

                  <Button
                    type="submit"
                    variant="primary"
                    size="lg"
                    loading={loading}
                    className="w-full group"
                    icon={<Shield size={20} />}
                  >
                    <span className="flex items-center space-x-2">
                      <span>Войти в Vira.ID</span>
                      <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform duration-300" />
                    </span>
                  </Button>
                </>
              ) : (
                <>
                  <div className="text-center space-y-6">
                    <div className="relative">
                      <div className="w-20 h-20 mx-auto rounded-full bg-gradient-to-br from-purple-400 to-pink-400 flex items-center justify-center shadow-lg shadow-purple-500/25">
                        <Code className="w-10 h-10 text-white" />
                      </div>
                      <div className="absolute -top-2 -right-2 w-8 h-8 bg-green-500 rounded-full flex items-center justify-center">
                        <CheckCircle className="w-5 h-5 text-white" />
                      </div>
                    </div>
                    
                    <div>
                      <h3 className="text-2xl font-bold text-white mb-3">
                        Добро пожаловать в Vira.Dev!
                      </h3>
                      <p className="text-white/70 leading-relaxed">
                        Это ваш первый вход в нашу образовательную платформу. 
                        Мы автоматически создадим для вас персональный профиль разработчика.
                      </p>
                    </div>

                    {/* Features preview */}
                    <div className="grid grid-cols-2 gap-4 text-left">
                      {[
                        { icon: '🎯', title: 'Персональные курсы', desc: 'Адаптированные под ваш уровень' },
                        { icon: '🚀', title: 'Практические проекты', desc: 'Реальные задачи из индустрии' },
                        { icon: '👥', title: 'Сообщество', desc: 'Общение с разработчиками' },
                        { icon: '🏆', title: 'Достижения', desc: 'Отслеживание прогресса' }
                      ].map((feature, index) => (
                        <div key={index} className="p-3 rounded-xl bg-white/5 hover:bg-white/10 transition-colors duration-300">
                          <div className="text-2xl mb-2">{feature.icon}</div>
                          <h4 className="text-white font-medium text-sm mb-1">{feature.title}</h4>
                          <p className="text-white/60 text-xs">{feature.desc}</p>
                        </div>
                      ))}
                    </div>
                  </div>

                  <Button
                    type="submit"
                    variant="primary"
                    size="lg"
                    loading={loading}
                    className="w-full group"
                    icon={<Code size={20} />}
                  >
                    <span className="flex items-center space-x-2">
                      <span>Создать профиль Vira.Dev</span>
                      <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform duration-300" />
                    </span>
                  </Button>
                </>
              )}
            </form>

            {/* Additional Options */}
            {step === 'vira-id' && (
              <div className="mt-8 pt-6 border-t border-white/10">
                <div className="text-center space-y-4">
                  <button 
                    onClick={() => navigate('/register')}
                    className="text-purple-400 hover:text-purple-300 transition-colors duration-300 font-medium"
                  >
                    Нет аккаунта? Создать новый
                  </button>
                  
                  <div className="flex items-center space-x-4 text-white/50 text-sm">
                    <div className="flex-1 h-px bg-white/10" />
                    <span>или войти через</span>
                    <div className="flex-1 h-px bg-white/10" />
                  </div>
                  
                  <div className="flex space-x-3">
                    <button className="flex-1 p-3 rounded-xl bg-white/5 hover:bg-white/10 transition-colors duration-300 text-white/70 hover:text-white">
                      <span className="text-xl">🔗</span>
                    </button>
                    <button className="flex-1 p-3 rounded-xl bg-white/5 hover:bg-white/10 transition-colors duration-300 text-white/70 hover:text-white">
                      <span className="text-xl">📱</span>
                    </button>
                  </div>
                </div>
              </div>
            )}
          </GlassCard>
        )}

        {/* Enhanced Progress Indicator */}
        {step !== 'transition' && (
          <div className="flex justify-center mt-8 space-x-3">
            <div className={`h-2 rounded-full transition-all duration-700 ${
              step === 'vira-id' ? 'bg-gradient-to-r from-blue-400 to-purple-400 w-8 shadow-lg shadow-blue-400/30' : 'bg-white/20 w-2'
            }`} />
            <div className={`h-2 rounded-full transition-all duration-700 ${
              step === 'vira-dev' ? 'bg-gradient-to-r from-purple-400 to-pink-400 w-8 shadow-lg shadow-purple-400/30' : 'bg-white/20 w-2'
            }`} />
          </div>
        )}

        {/* Step Labels */}
        {step !== 'transition' && (
          <div className="flex justify-between mt-4 text-sm">
            <span className={`transition-colors duration-300 ${
              step === 'vira-id' ? 'text-blue-400 font-medium' : 'text-white/50'
            }`}>
              Vira.ID
            </span>
            <span className={`transition-colors duration-300 ${
              step === 'vira-dev' ? 'text-purple-400 font-medium' : 'text-white/50'
            }`}>
              Vira.Dev
            </span>
          </div>
        )}
      </div>
    </div>
  );
}