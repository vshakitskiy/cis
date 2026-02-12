import { useAuth } from '@/context/AuthContext';
import { useNavigate, Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';

interface HeaderProps {
  title: string;
}

export function Header({ title }: HeaderProps) {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <header className="h-14 border-b border-border bg-white flex items-center px-6">
      <Link to="/" className="text-lg font-bold text-foreground hover:opacity-80 transition-opacity">
        CIS
      </Link>
      <Separator orientation="vertical" className="mx-4 h-5" />
      <h1 className="text-sm font-medium text-muted-foreground">{title}</h1>
      <div className="ml-auto flex items-center gap-3">
        {user && (
          <span className="text-sm text-muted-foreground capitalize">{user.role}</span>
        )}
        <Button variant="ghost" size="sm" onClick={handleLogout}>
          Logout
        </Button>
      </div>
    </header>
  );
}
