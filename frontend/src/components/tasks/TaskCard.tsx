import { Link } from 'react-router-dom';
import type { Task } from '@/types';
import { StatusBadge } from './StatusBadge';
import { PriorityBadge } from './PriorityBadge';
import { formatDate } from '@/lib/utils';
import { cn } from '@/lib/utils';

interface TaskCardProps {
  task: Task;
  projectId: number;
  draggable?: boolean;
}

export function TaskCard({ task, projectId, draggable }: TaskCardProps) {
  const handleDragStart = (e: React.DragEvent) => {
    e.dataTransfer.setData('text/plain', String(task.id));
    e.dataTransfer.effectAllowed = 'move';
  };

  return (
    <Link
      to={`/projects/${projectId}/tasks/${task.id}`}
      draggable={draggable}
      onDragStart={draggable ? handleDragStart : undefined}
      className={cn(
        'block bg-white rounded-md border border-gray-200 p-3 hover:shadow-sm transition-shadow',
        draggable && 'cursor-grab active:cursor-grabbing active:opacity-70',
      )}
    >
      <h4 className="text-sm font-medium text-gray-900 mb-2">{task.title}</h4>
      <div className="flex items-center gap-2 flex-wrap">
        <PriorityBadge priority={task.priority} />
        <StatusBadge status={task.status} />
        {task.deadline && (
          <span className="text-xs text-gray-400">Due {formatDate(task.deadline)}</span>
        )}
      </div>
    </Link>
  );
}
