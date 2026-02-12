import type { TaskPriority } from '@/types';
import { Badge } from '@/components/ui/badge';
import { PRIORITY_LABELS, PRIORITY_COLORS } from '@/lib/constants';

export function PriorityBadge({ priority }: { priority: TaskPriority }) {
  return <Badge variant="outline" className={PRIORITY_COLORS[priority]}>{PRIORITY_LABELS[priority]}</Badge>;
}
