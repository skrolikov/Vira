import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { 
  ArrowRight, 
  Code, 
  Users, 
  Trophy, 
  BookOpen, 
  Target, 
  Play,
  Sparkles,
  Rocket,
  Brain,
  Layers,
  GitBranch
} from 'lucide-react';
import { Button } from '../components/button';
import { GlassCard } from '../components/glassCard';

interface AnimatedCounterProps {
  end: number;
  duration?: number;
  suffix?: string;
  prefix?: string;
}

function AnimatedCounter({ end, duration = 2000, suffix = '', prefix = '' }: AnimatedCounterProps) {
  const [count, setCount] = useState(0);

  useEffect(() => {
    let startTime: number;
    let animationFrame: number;

    const animate = (timestamp: number) => {
      if (!startTime) startTime = timestamp;
      const progress = Math.min((timestamp - startTime) / duration, 1);
      
      setCount(Math.floor(progress * end));
      
      if (progress < 1) {
        animationFrame = requestAnimationFrame(animate);
      }
    };

    animationFrame = requestAnimationFrame(animate);
    return () => cancelAnimationFrame(animationFrame);
  }, [end, duration]);

  return <span>{prefix}{count.toLocaleString()}{suffix}</span>;
}

export function Home() {
  const [hoveredCard, setHoveredCard] = useState<string | null>(null);

  const features = [
    {
      id: 'interactive',
      title: 'Интерактивное обучение',
      description: 'Практические задания с мгновенной обратной связью',
      icon: <Play className="w-8 h-8" />,
      gradient: 'from-blue-500 to-cyan-500',
      stats: { value: 500, suffix: '+', label: 'Интерактивных уроков' }
    },
    {
      id: 'community',
      title: 'Сообщество разработчиков',
      description: 'Общайтесь, делитесь опытом и растите вместе',
      icon: <Users className="w-8 h-8" />,
      gradient: 'from-purple-500 to-pink-500',
      stats: { value: 25000, suffix: '+', label: 'Активных участников' }
    },
    {
      id: 'projects',
      title: 'Реальные проекты',
      description: 'Создавайте портфолио на практических задачах',
      icon: <Code className="w-8 h-8" />,
      gradient: 'from-green-500 to-emerald-500',
      stats: { value: 150, suffix: '+', label: 'Проектов в портфолио' }
    },
    {
      id: 'mentorship',
      title: 'Менторство',
      description: 'Персональная поддержка от опытных разработчиков',
      icon: <Target className="w-8 h-8" />,
      gradient: 'from-orange-500 to-red-500',
      stats: { value: 100, suffix: '+', label: 'Экспертов-менторов' }
    }
  ];

  const technologies = [
    { name: 'React', icon: '⚛️', level: 95 },
    { name: 'TypeScript', icon: '📘', level: 90 },
    { name: 'Node.js', icon: '🟢', level: 88 },
    { name: 'Python', icon: '🐍', level: 85 },
    { name: 'Docker', icon: '🐳', level: 82 },
    { name: 'AWS', icon: '☁️', level: 78 }
  ];

  const achievements = [
    { icon: '🏆', title: 'Лучшая EdTech платформа 2024', org: 'TechCrunch' },
    { icon: '⭐', title: '4.9/5 рейтинг студентов', org: 'Trustpilot' },
    { icon: '🚀', title: '95% трудоустройство выпускников', org: 'Career Report' }
  ];

  return (
    <div className="min-h-screen relative overflow-hidden">
      {/* Enhanced animated background */}
      <div className="absolute inset-0">

        {/* Floating code symbols */}
        {Array.from({ length: 20 }).map((_, i) => (
          <div
            key={i}
            className="absolute text-white/10 font-mono text-2xl animate-float"
            style={{
              left: `${Math.random() * 100}%`,
              top: `${Math.random() * 100}%`,
              animationDelay: `${Math.random() * 20}s`,
              animationDuration: `${15 + Math.random() * 10}s`,
            }}
          >
            {['<>', '{}', '[]', '()', '/>', '&&', '||', '=>'][Math.floor(Math.random() * 8)]}
          </div>
        ))}
      </div>

      <div className="relative z-10 p-6">
        <div className="max-w-7xl mx-auto">
          {/* Hero Section */}
          <div className="text-center mb-16 pt-20">
            <div className="flex items-center justify-center space-x-3 mb-8 animate-fade-in-up">
              <Code className="w-12 h-12 text-purple-400" />
              <h1 className="text-5xl md:text-7xl font-bold bg-gradient-to-r from-purple-400 via-pink-400 to-cyan-400 bg-clip-text text-transparent">
                Vira.Dev
              </h1>
              <Sparkles className="w-8 h-8 text-yellow-400 animate-pulse" />
            </div>
            
            <h2 className="text-3xl md:text-5xl font-bold text-white mb-6 animate-fade-in-up delay-200">
              Будущее разработки
              <span className="block bg-gradient-to-r from-purple-400 to-pink-400 bg-clip-text text-transparent">
                начинается здесь
              </span>
            </h2>
            
            <p className="text-xl text-white/80 mb-8 max-w-3xl mx-auto leading-relaxed animate-fade-in-up delay-300">
              Революционная платформа для изучения программирования с ИИ-ментором, 
              интерактивными проектами и сообществом мирового уровня
            </p>

            <div className="flex flex-col sm:flex-row gap-4 justify-center mb-12 animate-fade-in-up delay-500">
              <Link to="/register">
                <Button
                  variant="primary"
                  size="lg"
                  icon={<Rocket size={20} />}
                  className="group relative overflow-hidden"
                >
                  <span className="relative z-10">Начать обучение</span>
                  <div className="absolute inset-0 bg-gradient-to-r from-white/0 via-white/20 to-white/0 -translate-x-full group-hover:translate-x-full transition-transform duration-1000" />
                </Button>
              </Link>
              
              <Link to="/auth">
                <Button
                  variant="secondary"
                  size="lg"
                  icon={<ArrowRight size={20} />}
                  className="group"
                >
                  <span>Войти в аккаунт</span>
                </Button>
              </Link>
            </div>

            {/* Live stats */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-6 max-w-4xl mx-auto animate-fade-in-up delay-700">
              {[
                { value: 50000, suffix: '+', label: 'Студентов' },
                { value: 1200, suffix: '+', label: 'Курсов' },
                { value: 95, suffix: '%', label: 'Трудоустройство' },
                { value: 24, suffix: '/7', label: 'Поддержка' }
              ].map((stat, index) => (
                <div key={index} className="text-center">
                  <div className="text-3xl font-bold text-white mb-2">
                    <AnimatedCounter end={stat.value} suffix={stat.suffix} />
                  </div>
                  <div className="text-white/60 text-sm">{stat.label}</div>
                </div>
              ))}
            </div>
          </div>

          {/* Bento Grid Layout */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-16">
            {/* Large feature card */}
            <div className="lg:col-span-2 lg:row-span-2">
              <GlassCard className="p-8 h-full relative group overflow-hidden">
                <div className="absolute inset-0 bg-gradient-to-br from-purple-500/10 to-pink-500/10 opacity-0 group-hover:opacity-100 transition-opacity duration-500" />
                <div className="relative z-10">
                  <div className="flex items-center space-x-3 mb-6">
                    <div className="p-3 rounded-2xl bg-gradient-to-br from-purple-500/20 to-pink-500/20">
                      <Brain className="w-8 h-8 text-purple-400" />
                    </div>
                    <h3 className="text-2xl font-bold text-white">ИИ-Ментор</h3>
                  </div>
                  
                  <p className="text-white/70 text-lg mb-6 leading-relaxed">
                    Персональный искусственный интеллект анализирует ваш прогресс, 
                    адаптирует программу обучения и предоставляет мгновенную обратную связь
                  </p>
                  
                  <div className="space-y-4 mb-8">
                    {[
                      'Анализ кода в реальном времени',
                      'Персонализированные рекомендации',
                      'Автоматическая проверка заданий',
                      'Прогнозирование сложностей'
                    ].map((feature, index) => (
                      <div key={index} className="flex items-center space-x-3">
                        <div className="w-2 h-2 bg-purple-400 rounded-full" />
                        <span className="text-white/80">{feature}</span>
                      </div>
                    ))}
                  </div>
                  
                  <Button variant="primary" icon={<Sparkles size={18} />}>
                    Попробовать ИИ-ментора
                  </Button>
                </div>
              </GlassCard>
            </div>

            {/* Technology stack */}
            <div className="lg:col-span-1">
              <GlassCard className="p-6 h-full">
                <div className="flex items-center space-x-3 mb-6">
                  <Layers className="w-6 h-6 text-cyan-400" />
                  <h3 className="text-xl font-bold text-white">Технологии</h3>
                </div>
                
                <div className="space-y-4">
                  {technologies.slice(0, 4).map((tech, index) => (
                    <div key={index} className="group">
                      <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center space-x-2">
                          <span className="text-xl">{tech.icon}</span>
                          <span className="text-white font-medium">{tech.name}</span>
                        </div>
                        <span className="text-white/60 text-sm">{tech.level}%</span>
                      </div>
                      <div className="w-full bg-white/10 rounded-full h-2">
                        <div 
                          className="bg-gradient-to-r from-cyan-500 to-blue-500 h-2 rounded-full transition-all duration-1000 group-hover:from-purple-500 group-hover:to-pink-500"
                          style={{ width: `${tech.level}%` }}
                        />
                      </div>
                    </div>
                  ))}
                </div>
              </GlassCard>
            </div>

            {/* Live coding preview */}
            <div className="lg:col-span-1">
              <GlassCard className="p-6 h-full bg-black/20">
                <div className="flex items-center space-x-3 mb-4">
                  <div className="flex space-x-1">
                    <div className="w-3 h-3 bg-red-400 rounded-full" />
                    <div className="w-3 h-3 bg-yellow-400 rounded-full" />
                    <div className="w-3 h-3 bg-green-400 rounded-full" />
                  </div>
                  <span className="text-white/60 text-sm">live-coding.tsx</span>
                </div>
                
                <div className="font-mono text-sm space-y-2">
                  <div className="text-purple-400">
  <span className="text-white">const </span>
  <span className="text-cyan-400">learn</span>
  <span className="text-white"> = () =&gt; {'{'}</span>
</div>
                  <div className="text-white/80 ml-4">return (</div>
                  <div className="text-green-400 ml-8">&lt;Future</div>
                  <div className="text-blue-400 ml-12">career={`{`}<span className="text-yellow-400">"amazing"</span>{`}`}</div>
                  <div className="text-blue-400 ml-12">skills={`{`}<span className="text-yellow-400">"unlimited"</span>{`}`}</div>
                  <div className="text-green-400 ml-8">/&gt;</div>
                  <div className="text-white/80 ml-4">);</div>
                  <div className="text-purple-400">{`}`};</div>
                </div>
                
                <div className="mt-4 flex items-center space-x-2">
                  <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse" />
                  <span className="text-green-400 text-xs">Компиляция успешна</span>
                </div>
              </GlassCard>
            </div>

            {/* Features grid */}
            {features.map((feature, index) => (
              <div key={feature.id} className="lg:col-span-1">
                <GlassCard 
                  className={`p-6 h-full cursor-pointer transition-all duration-500 hover:scale-105 ${
                    hoveredCard === feature.id ? 'shadow-2xl shadow-purple-500/20' : ''
                  }`}
                  onMouseEnter={() => setHoveredCard(feature.id)}
                  onMouseLeave={() => setHoveredCard(null)}
                >
                  <div className={`p-3 rounded-2xl bg-gradient-to-br ${feature.gradient}/20 mb-4 w-fit`}>
                    <div className={`text-transparent bg-gradient-to-r ${feature.gradient} bg-clip-text`}>
                      {feature.icon}
                    </div>
                  </div>
                  
                  <h3 className="text-lg font-bold text-white mb-3">{feature.title}</h3>
                  <p className="text-white/70 text-sm mb-4 leading-relaxed">{feature.description}</p>
                  
                  <div className="text-center">
                    <div className="text-2xl font-bold text-white mb-1">
                      <AnimatedCounter 
                        end={feature.stats.value} 
                        suffix={feature.stats.suffix}
                        duration={2000 + index * 200}
                      />
                    </div>
                    <div className="text-white/60 text-xs">{feature.stats.label}</div>
                  </div>
                </GlassCard>
              </div>
            ))}

            {/* Achievements */}
            <div className="lg:col-span-2">
              <GlassCard className="p-6 h-full">
                <div className="flex items-center space-x-3 mb-6">
                  <Trophy className="w-6 h-6 text-yellow-400" />
                  <h3 className="text-xl font-bold text-white">Достижения</h3>
                </div>
                
                <div className="grid gap-4">
                  {achievements.map((achievement, index) => (
                    <div key={index} className="flex items-center space-x-4 p-4 rounded-xl bg-white/5 hover:bg-white/10 transition-colors duration-300">
                      <span className="text-3xl">{achievement.icon}</span>
                      <div>
                        <h4 className="text-white font-medium">{achievement.title}</h4>
                        <p className="text-white/60 text-sm">{achievement.org}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </GlassCard>
            </div>

            {/* Learning path preview */}
            <div className="lg:col-span-2">
              <GlassCard className="p-6 h-full">
                <div className="flex items-center space-x-3 mb-6">
                  <GitBranch className="w-6 h-6 text-green-400" />
                  <h3 className="text-xl font-bold text-white">Путь обучения</h3>
                </div>
                
                <div className="space-y-4">
                  {[
                    { step: 1, title: 'Основы программирования', status: 'completed', progress: 100 },
                    { step: 2, title: 'Frontend разработка', status: 'current', progress: 65 },
                    { step: 3, title: 'Backend и базы данных', status: 'locked', progress: 0 },
                    { step: 4, title: 'DevOps и деплой', status: 'locked', progress: 0 }
                  ].map((item, index) => (
                    <div key={index} className="flex items-center space-x-4">
                      <div className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold ${
                        item.status === 'completed' ? 'bg-green-500 text-white' :
                        item.status === 'current' ? 'bg-purple-500 text-white' :
                        'bg-white/20 text-white/60'
                      }`}>
                        {item.status === 'completed' ? '✓' : item.step}
                      </div>
                      
                      <div className="flex-1">
                        <div className="flex justify-between items-center mb-1">
                          <span className={`font-medium ${
                            item.status === 'locked' ? 'text-white/60' : 'text-white'
                          }`}>
                            {item.title}
                          </span>
                          <span className="text-white/60 text-sm">{item.progress}%</span>
                        </div>
                        <div className="w-full bg-white/10 rounded-full h-2">
                          <div 
                            className={`h-2 rounded-full transition-all duration-1000 ${
                              item.status === 'completed' ? 'bg-green-500' :
                              item.status === 'current' ? 'bg-gradient-to-r from-purple-500 to-pink-500' :
                              'bg-white/20'
                            }`}
                            style={{ width: `${item.progress}%` }}
                          />
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </GlassCard>
            </div>
          </div>

          {/* CTA Section */}
          <div className="text-center">
            <GlassCard className="p-12 relative overflow-hidden group">
              <div className="absolute inset-0 bg-gradient-to-r from-purple-500/10 via-pink-500/10 to-cyan-500/10 opacity-0 group-hover:opacity-100 transition-opacity duration-1000" />
              <div className="relative z-10">
                <h2 className="text-4xl font-bold text-white mb-6">
                  Готовы изменить свою жизнь?
                </h2>
                <p className="text-xl text-white/80 mb-8 max-w-2xl mx-auto">
                  Присоединяйтесь к тысячам разработчиков, которые уже строят свое будущее с Vira.Dev
                </p>
                
                <div className="flex flex-col sm:flex-row gap-4 justify-center">
                  <Link to="/register">
                    <Button
                      variant="primary"
                      size="lg"
                      icon={<Rocket size={20} />}
                      className="group relative overflow-hidden"
                    >
                      <span className="relative z-10">Начать бесплатно</span>
                      <div className="absolute inset-0 bg-gradient-to-r from-white/0 via-white/20 to-white/0 -translate-x-full group-hover:translate-x-full transition-transform duration-1000" />
                    </Button>
                  </Link>
                  
                  <Button
                    variant="secondary"
                    size="lg"
                    icon={<BookOpen size={20} />}
                  >
                    Смотреть демо
                  </Button>
                </div>
              </div>
            </GlassCard>
          </div>
        </div>
      </div>
    </div>
  );
}