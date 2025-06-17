interface GlassCardProps {
    children: React.ReactNode;
    className?: string;
    animate?: boolean;
    onMouseEnter?: React.MouseEventHandler<HTMLDivElement>;
    onMouseLeave?: React.MouseEventHandler<HTMLDivElement>;
}

export function GlassCard({ children, className = '', animate = true }: GlassCardProps) {
    return (
        <div className={`
      glass-card backdrop-blur-2xl border border-white/10 dark:border-white/5
      ${animate ? 'animate-fade-in-up' : ''} 
      ${className}
    `}>
            {children}
        </div>
    );
}