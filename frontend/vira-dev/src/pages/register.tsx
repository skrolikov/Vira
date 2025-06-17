import { useState } from 'react';
import { 
  Mail, 
  Lock, 
  User, 
  ArrowRight, 
  Code, 
  CheckCircle, 
  Eye,
  EyeOff,
  Zap,
  Target,
  BookOpen,
} from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { InputField } from '../components/inputField';
import { ServiceBranding } from '../components/serviceBranding';
import { GlassCard } from '../components/glassCard';
import { Button } from '../components/button';

type RegisterStep = 1 | 2 | 3;

interface FormData {
  email: string;
  fullName: string;
  password: string;
  confirmPassword: string;
  username: string;
  interests: string[];
  experience: string;
  goals: string[];
}

export function Register() {
  const [step, setStep] = useState<RegisterStep>(1);
  const [loading, setLoading] = useState(false);
  const [isTransitioning, setIsTransitioning] = useState(false);
  const navigate = useNavigate();
  
  const [formData, setFormData] = useState<FormData>({
    email: '',
    fullName: '',
    password: '',
    confirmPassword: '',
    username: '',
    interests: [],
    experience: 'beginner',
    goals: []
  });

  const [passwordStrength, setPasswordStrength] = useState(0);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const calculatePasswordStrength = (password: string) => {
    let strength = 0;
    if (password.length >= 8) strength += 25;
    if (/[A-Z]/.test(password)) strength += 25;
    if (/[0-9]/.test(password)) strength += 25;
    if (/[^A-Za-z0-9]/.test(password)) strength += 25;
    return strength;
  };

  const handlePasswordChange = (password: string) => {
    setFormData({ ...formData, password });
    setPasswordStrength(calculatePasswordStrength(password));
  };

  const handleNext = async () => {
    setLoading(true);
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1200));
    setLoading(false);
    
    if (step < 3) {
      setIsTransitioning(true);
      setTimeout(() => {
        setStep((prev) => (prev + 1) as RegisterStep);
        setIsTransitioning(false);
      }, 600);
    } else {
      // Registration complete
      navigate('/auth');
    }
  };

  const handleInterestToggle = (interest: string) => {
    setFormData(prev => ({
      ...prev,
      interests: prev.interests.includes(interest)
        ? prev.interests.filter(i => i !== interest)
        : [...prev.interests, interest]
    }));
  };

  const handleGoalToggle = (goal: string) => {
    setFormData(prev => ({
      ...prev,
      goals: prev.goals.includes(goal)
        ? prev.goals.filter(g => g !== goal)
        : [...prev.goals, goal]
    }));
  };

  const getStepContent = () => {
    switch (step) {
      case 1:
        return {
          title: 'Создание аккаунта',
          subtitle: 'Начнем с основной информации для вашего Vira.ID',
          service: 'vira-id' as const,
          progress: 33,
          fields: (
            <div className="space-y-5">
              <InputField
                type="email"
                placeholder="Введите ваш email"
                value={formData.email}
                onChange={(value) => setFormData({ ...formData, email: value })}
                icon={<Mail size={20} />}
                required
              />
              <InputField
                type="text"
                placeholder="Полное имя"
                value={formData.fullName}
                onChange={(value) => setFormData({ ...formData, fullName: value })}
                icon={<User size={20} />}
                required
              />
              
              {/* Email verification hint */}
              <div className="p-4 rounded-xl bg-blue-500/10 border border-blue-400/20">
                <div className="flex items-start space-x-3">
                  <Mail className="w-5 h-5 text-blue-400 mt-0.5" />
                  <div>
                    <p className="text-blue-400 font-medium text-sm">Подтверждение email</p>
                    <p className="text-white/70 text-xs mt-1">
                      На указанный адрес будет отправлено письмо для подтверждения
                    </p>
                  </div>
                </div>
              </div>
            </div>
          )
        };
      
      case 2:
        return {
          title: 'Защита аккаунта',
          subtitle: 'Создайте надежный пароль для безопасности',
          service: 'vira-id' as const,
          progress: 66,
          fields: (
            <div className="space-y-5">
              <div>
                <div className="relative">
                  <InputField
                    type={showPassword ? "text" : "password"}
                    placeholder="Создайте пароль"
                    value={formData.password}
                    onChange={handlePasswordChange}
                    icon={<Lock size={20} />}
                    required
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-4 top-1/2 -translate-y-1/2 text-white/60 hover:text-purple-400 transition-colors duration-300"
                  >
                    {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                  </button>
                </div>
                
                {/* Password strength indicator */}
                {formData.password && (
                  <div className="mt-3">
                    <div className="flex justify-between text-xs text-white/60 mb-2">
                      <span>Надежность пароля</span>
                      <span>{passwordStrength}%</span>
                    </div>
                    <div className="w-full bg-white/10 rounded-full h-2">
                      <div 
                        className={`h-2 rounded-full transition-all duration-500 ${
                          passwordStrength < 50 ? 'bg-red-400' :
                          passwordStrength < 75 ? 'bg-yellow-400' :
                          'bg-green-400'
                        }`}
                        style={{ width: `${passwordStrength}%` }}
                      />
                    </div>
                  </div>
                )}
              </div>
              
              <div className="relative">
                <InputField
                  type={showConfirmPassword ? "text" : "password"}
                  placeholder="Подтвердите пароль"
                  value={formData.confirmPassword}
                  onChange={(value) => setFormData({ ...formData, confirmPassword: value })}
                  icon={<Lock size={20} />}
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  className="absolute right-4 top-1/2 -translate-y-1/2 text-white/60 hover:text-purple-400 transition-colors duration-300"
                >
                  {showConfirmPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                </button>
              </div>
              
              {/* Password requirements */}
              <div className="p-4 rounded-xl bg-white/5 border border-white/10">
                <p className="text-white/70 text-sm mb-3">Пароль должен содержать:</p>
                <div className="grid grid-cols-2 gap-2 text-xs">
                  {[
                    { text: 'Минимум 8 символов', check: formData.password.length >= 8 },
                    { text: 'Заглавные буквы', check: /[A-Z]/.test(formData.password) },
                    { text: 'Строчные буквы', check: /[a-z]/.test(formData.password) },
                    { text: 'Цифры', check: /[0-9]/.test(formData.password) },
                  ].map((req, index) => (
                    <div key={index} className={`flex items-center space-x-2 ${req.check ? 'text-green-400' : 'text-white/50'}`}>
                      <CheckCircle className={`w-3 h-3 ${req.check ? 'opacity-100' : 'opacity-30'}`} />
                      <span>{req.text}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )
        };
      
      case 3:
        return {
          title: 'Персонализация профиля',
          subtitle: 'Настройте ваш учебный профиль в Vira.Dev',
          service: 'vira-dev' as const,
          progress: 100,
          fields: (
            <div className="space-y-6">
              <InputField
                type="text"
                placeholder="Выберите имя пользователя"
                value={formData.username}
                onChange={(value) => setFormData({ ...formData, username: value })}
                icon={<Code size={20} />}
                required
              />
              
              {/* Experience level */}
              <div>
                <h4 className="text-white font-medium mb-4 flex items-center space-x-2">
                  <Target className="w-5 h-5 text-purple-400" />
                  <span>Ваш уровень опыта</span>
                </h4>
                <div className="grid grid-cols-1 gap-3">
                  {[
                    { id: 'beginner', label: 'Новичок', desc: 'Только начинаю изучать программирование', icon: '🌱' },
                    { id: 'intermediate', label: 'Средний', desc: 'Есть базовые знания, хочу углубиться', icon: '🚀' },
                    { id: 'advanced', label: 'Продвинутый', desc: 'Опытный разработчик, изучаю новые технологии', icon: '⚡' }
                  ].map((level) => (
                    <label key={level.id} className="cursor-pointer group">
                      <input
                        type="radio"
                        name="experience"
                        value={level.id}
                        checked={formData.experience === level.id}
                        onChange={(e) => setFormData({ ...formData, experience: e.target.value })}
                        className="sr-only"
                      />
                      <div className={`p-4 rounded-xl border-2 transition-all duration-300 ${
                        formData.experience === level.id
                          ? 'border-purple-400 bg-purple-500/10'
                          : 'border-white/20 bg-white/5 group-hover:border-purple-400/50'
                      }`}>
                        <div className="flex items-start space-x-3">
                          <span className="text-2xl">{level.icon}</span>
                          <div>
                            <h5 className="text-white font-medium">{level.label}</h5>
                            <p className="text-white/60 text-sm">{level.desc}</p>
                          </div>
                        </div>
                      </div>
                    </label>
                  ))}
                </div>
              </div>

              {/* Learning interests */}
              <div>
                <h4 className="text-white font-medium mb-4 flex items-center space-x-2">
                  <BookOpen className="w-5 h-5 text-purple-400" />
                  <span>Направления изучения</span>
                </h4>
                <div className="grid grid-cols-2 gap-3">
                  {[
                    { id: 'frontend', label: 'Frontend', icon: '🎨' },
                    { id: 'backend', label: 'Backend', icon: '⚙️' },
                    { id: 'mobile', label: 'Mobile', icon: '📱' },
                    { id: 'devops', label: 'DevOps', icon: '🔧' },
                    { id: 'ai', label: 'AI/ML', icon: '🤖' },
                    { id: 'blockchain', label: 'Blockchain', icon: '⛓️' }
                  ].map((interest) => (
                    <label key={interest.id} className="cursor-pointer group">
                      <input
                        type="checkbox"
                        checked={formData.interests.includes(interest.id)}
                        onChange={() => handleInterestToggle(interest.id)}
                        className="sr-only"
                      />
                      <div className={`p-3 rounded-xl border-2 transition-all duration-300 text-center ${
                        formData.interests.includes(interest.id)
                          ? 'border-purple-400 bg-purple-500/10'
                          : 'border-white/20 bg-white/5 group-hover:border-purple-400/50'
                      }`}>
                        <div className="text-2xl mb-2">{interest.icon}</div>
                        <span className="text-white text-sm font-medium">{interest.label}</span>
                      </div>
                    </label>
                  ))}
                </div>
              </div>

              {/* Learning goals */}
              <div>
                <h4 className="text-white font-medium mb-4 flex items-center space-x-2">
                  <Zap className="w-5 h-5 text-purple-400" />
                  <span>Ваши цели</span>
                </h4>
                <div className="space-y-2">
                  {[
                    'Получить работу разработчика',
                    'Создать собственный проект',
                    'Повысить квалификацию',
                    'Изучить новые технологии',
                    'Подготовиться к собеседованиям'
                  ].map((goal) => (
                    <label key={goal} className="flex items-center space-x-3 cursor-pointer group p-2 rounded-lg hover:bg-white/5 transition-colors duration-300">
                      <input
                        type="checkbox"
                        checked={formData.goals.includes(goal)}
                        onChange={() => handleGoalToggle(goal)}
                        className="sr-only"
                      />
                      <div className={`w-5 h-5 rounded-lg border-2 transition-all duration-300 flex items-center justify-center ${
                        formData.goals.includes(goal)
                          ? 'border-purple-400 bg-purple-500'
                          : 'border-white/30 group-hover:border-purple-400'
                      }`}>
                        {formData.goals.includes(goal) && (
                          <CheckCircle className="w-3 h-3 text-white" />
                        )}
                      </div>
                      <span className="text-white/80 group-hover:text-white transition-colors duration-300 text-sm">
                        {goal}
                      </span>
                    </label>
                  ))}
                </div>
              </div>
            </div>
          )
        };
    }
  };

  const stepContent = getStepContent();

  return (
    <div className="min-h-screen flex items-center justify-center p-4 relative overflow-hidden">
      {/* Enhanced background effects */}
      <div className="absolute inset-0 bg-gradient-to-br from-purple-900/20 via-transparent to-pink-900/20" />
      
      <div className="w-full max-w-lg relative z-10">
        {/* Service Branding */}
        <div className="text-center mb-8">
          <div className={`transition-all duration-700 ease-out ${
            isTransitioning ? 'opacity-0 transform scale-95' : 'opacity-100 transform scale-100'
          }`}>
            <ServiceBranding
              service={stepContent.service} 
              className="justify-center mb-6"
            />
            <div>
              <h1 className="text-3xl font-bold text-white mb-3">
                {stepContent.title}
              </h1>
              <p className="text-white/70 leading-relaxed">
                {stepContent.subtitle}
              </p>
            </div>
          </div>
        </div>

        {/* Progress Bar */}
        <div className="mb-8">
          <div className="flex justify-between text-sm text-white/60 mb-2">
            <span>Шаг {step} из 3</span>
            <span>{stepContent.progress}%</span>
          </div>
          <div className="w-full bg-white/10 rounded-full h-2">
            <div 
              className="bg-gradient-to-r from-purple-500 to-pink-500 h-2 rounded-full transition-all duration-1000 ease-out"
              style={{ width: `${stepContent.progress}%` }}
            />
          </div>
        </div>

        {/* Registration Form */}
        <GlassCard className={`p-8 transition-all duration-700 ease-out ${
          isTransitioning ? 'opacity-0 transform translate-x-8 scale-95' : 'opacity-100 transform translate-x-0 scale-100'
        }`}>
          <form onSubmit={(e) => {
            e.preventDefault();
            handleNext();
          }}>
            
            {stepContent.fields}

            <div className="mt-8">
              <Button
                type="submit"
                variant="primary"
                size="lg"
                loading={loading}
                className="w-full group"
                icon={step < 3 ? <ArrowRight size={20} /> : <CheckCircle size={20} />}
              >
                <span className="flex items-center space-x-2">
                  <span>{step < 3 ? 'Продолжить' : 'Завершить регистрацию'}</span>
                  {step < 3 && (
                    <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform duration-300" />
                  )}
                </span>
              </Button>
            </div>
          </form>

          {/* Additional Options */}
          <div className="mt-6 pt-6 border-t border-white/10">
            <div className="text-center">
              <button 
                onClick={() => navigate('/auth')}
                className="text-purple-400 hover:text-purple-300 transition-colors duration-300 font-medium"
              >
                Уже есть аккаунт? Войти
              </button>
            </div>
          </div>
        </GlassCard>

        {/* Enhanced Step Indicators */}
        <div className="flex justify-center mt-8 space-x-4">
          {[1, 2, 3].map((stepNum) => (
            <div key={stepNum} className="flex flex-col items-center space-y-2">
              <div className={`w-10 h-10 rounded-full flex items-center justify-center transition-all duration-500 ${
                stepNum === step 
                  ? 'bg-gradient-to-r from-purple-500 to-pink-500 shadow-lg shadow-purple-500/30' 
                  : stepNum < step 
                    ? 'bg-green-500 shadow-lg shadow-green-500/30'
                    : 'bg-white/20'
              }`}>
                {stepNum < step ? (
                  <CheckCircle className="w-5 h-5 text-white" />
                ) : (
                  <span className="text-white font-medium">{stepNum}</span>
                )}
              </div>
              <span className={`text-xs transition-colors duration-300 ${
                stepNum <= step ? 'text-white' : 'text-white/50'
              }`}>
                {stepNum === 1 ? 'Основное' : stepNum === 2 ? 'Безопасность' : 'Профиль'}
              </span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}