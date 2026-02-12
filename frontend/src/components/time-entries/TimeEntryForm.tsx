import { useState, type FormEvent } from 'react';
import { FormField } from '@/components/ui/form-field';
import { Button } from '@/components/ui/button';

interface TimeEntryFormProps {
  onSubmit: (data: { minutes: number; description?: string; date: string }) => Promise<void>;
}

export function TimeEntryForm({ onSubmit }: TimeEntryFormProps) {
  const today = new Date().toISOString().slice(0, 10);
  const [minutes, setMinutes] = useState('');
  const [description, setDescription] = useState('');
  const [date, setDate] = useState(today);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    const mins = parseInt(minutes, 10);
    if (!mins || mins <= 0) return;
    setLoading(true);
    try {
      await onSubmit({ minutes: mins, description: description || undefined, date });
      setMinutes('');
      setDescription('');
      setDate(today);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-wrap gap-2 items-end">
      <FormField
        id="te-minutes"
        label="Minutes"
        type="number"
        min={1}
        value={minutes}
        onChange={(e) => setMinutes(e.target.value)}
        required
        className="w-24"
      />
      <FormField
        id="te-desc"
        label="Description"
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        className="flex-1 min-w-[150px]"
      />
      <FormField
        id="te-date"
        label="Date"
        type="date"
        value={date}
        onChange={(e) => setDate(e.target.value)}
        required
      />
      <Button type="submit" disabled={loading}>
        {loading ? '...' : 'Log Time'}
      </Button>
    </form>
  );
}
