import { useState } from 'react';
import type { Task, TaskStatus } from '@/types';
import { TaskCard } from './TaskCard';
import { STATUS_LABELS, STATUSES } from '@/lib/constants';
import { cn } from '@/lib/utils';

interface TaskBoardProps {
  tasks: Task[];
  projectId: number;
  onStatusChange?: (taskId: number, newStatus: TaskStatus) => void;
}

function Column({
  status,
  tasks,
  projectId,
  onDrop,
}: {
  status: TaskStatus;
  tasks: Task[];
  projectId: number;
  onDrop: (taskId: number, newStatus: TaskStatus) => void;
}) {
  const [dragOver, setDragOver] = useState(false);

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
    setDragOver(true);
  };

  const handleDragLeave = () => {
    setDragOver(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setDragOver(false);
    const taskId = Number(e.dataTransfer.getData('text/plain'));
    if (taskId) {
      onDrop(taskId, status);
    }
  };

  return (
    <div
      className={cn(
        'flex-1 min-w-[250px] rounded-lg p-3 transition-colors',
        dragOver ? 'bg-accent/50 ring-2 ring-ring/20' : 'bg-muted/40',
      )}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
    >
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-semibold text-gray-700">{STATUS_LABELS[status]}</h3>
        <span className="text-xs text-muted-foreground bg-background rounded-full px-2 py-0.5 border">
          {tasks.length}
        </span>
      </div>
      <div className="space-y-2">
        {tasks.map((t) => (
          <TaskCard key={t.id} task={t} projectId={projectId} draggable />
        ))}
        {tasks.length === 0 && (
          <p className="text-xs text-muted-foreground text-center py-8">No tasks</p>
        )}
      </div>
    </div>
  );
}

export function TaskBoard({ tasks, projectId, onStatusChange }: TaskBoardProps) {
  const grouped = STATUSES.reduce(
    (acc, s) => {
      acc[s] = tasks.filter((t) => t.status === s);
      return acc;
    },
    {} as Record<TaskStatus, Task[]>,
  );

  const handleDrop = (taskId: number, newStatus: TaskStatus) => {
    const task = tasks.find((t) => t.id === taskId);
    if (!task || task.status === newStatus) return;
    onStatusChange?.(taskId, newStatus);
  };

  return (
    <div className="flex gap-4 overflow-x-auto pb-2">
      {STATUSES.map((s) => (
        <Column key={s} status={s} tasks={grouped[s]} projectId={projectId} onDrop={handleDrop} />
      ))}
    </div>
  );
}
