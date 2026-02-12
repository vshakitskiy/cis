import { useState, type FormEvent } from 'react';
import type { Task, TaskStatus, TaskPriority } from '@/types';
import { FormField } from '@/components/ui/form-field';
import { FormSelect } from '@/components/ui/form-select';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { STATUSES, PRIORITIES, STATUS_LABELS, PRIORITY_LABELS } from '@/lib/constants';

interface TaskFormData {
  title: string;
  description?: string;
  assignee_id?: number | null;
  status: TaskStatus;
  priority: TaskPriority;
  deadline?: string;
}

interface TaskFormProps {
  initial?: Task;
  onSubmit: (data: TaskFormData) => Promise<void>;
  onCancel: () => void;
}

export function TaskForm({ initial, onSubmit, onCancel }: TaskFormProps) {
  const [title, setTitle] = useState(initial?.title ?? '');
  const [description, setDescription] = useState(initial?.description ?? '');
  const [status, setStatus] = useState<TaskStatus>(initial?.status ?? 'todo');
  const [priority, setPriority] = useState<TaskPriority>(initial?.priority ?? 'medium');
  const [deadline, setDeadline] = useState(initial?.deadline ?? '');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await onSubmit({
        title,
        description: description || undefined,
        status,
        priority,
        deadline: deadline || undefined,
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save');
    } finally {
      setLoading(false);
    }
  };

  const statusOptions = STATUSES.map((s) => ({ value: s, label: STATUS_LABELS[s] }));
  const priorityOptions = PRIORITIES.map((p) => ({ value: p, label: PRIORITY_LABELS[p] }));

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && <div className="rounded-md bg-red-50 p-3 text-sm text-red-600">{error}</div>}
      <FormField
        id="task-title"
        label="Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        required
      />
      <div className="space-y-1">
        <Label htmlFor="task-desc">Description</Label>
        <Textarea
          id="task-desc"
          rows={3}
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />
      </div>
      {initial && (
        <FormSelect
          id="task-status"
          label="Status"
          options={statusOptions}
          value={status}
          onValueChange={(v) => setStatus(v as TaskStatus)}
        />
      )}
      <FormSelect
        id="task-priority"
        label="Priority"
        options={priorityOptions}
        value={priority}
        onValueChange={(v) => setPriority(v as TaskPriority)}
      />
      <FormField
        id="task-deadline"
        label="Deadline"
        type="date"
        value={deadline}
        onChange={(e) => setDeadline(e.target.value)}
      />
      <div className="flex justify-end gap-2">
        <Button type="button" variant="secondary" onClick={onCancel}>
          Cancel
        </Button>
        <Button type="submit" disabled={loading}>
          {loading ? 'Saving...' : initial ? 'Update' : 'Create'}
        </Button>
      </div>
    </form>
  );
}
