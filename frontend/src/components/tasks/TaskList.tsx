import type { Task } from '@/types';
import { TaskCard } from './TaskCard';

interface TaskListProps {
  tasks: Task[];
  projectId: number;
}

export function TaskList({ tasks, projectId }: TaskListProps) {
  if (tasks.length === 0) {
    return (
      <p className="text-center py-8 text-gray-500">No tasks yet. Create one to get started.</p>
    );
  }

  return (
    <div className="space-y-2">
      {tasks.map((t) => (
        <TaskCard key={t.id} task={t} projectId={projectId} />
      ))}
    </div>
  );
}
