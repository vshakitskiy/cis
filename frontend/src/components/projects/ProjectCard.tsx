import { Link } from 'react-router-dom';
import type { Project } from '@/types';
import { formatDate } from '@/lib/utils';

interface ProjectCardProps {
  project: Project;
}

export function ProjectCard({ project }: ProjectCardProps) {
  return (
    <Link
      to={`/projects/${project.id}`}
      className="block bg-white rounded-lg border border-gray-200 p-5 hover:shadow-md transition-shadow"
    >
      <h3 className="font-semibold text-gray-900 mb-1">{project.name}</h3>
      {project.description && (
        <p className="text-sm text-gray-500 mb-3 line-clamp-2">{project.description}</p>
      )}
      <div className="flex items-center justify-between text-xs text-gray-400">
        <span>Created {formatDate(project.created_at)}</span>
        {project.deadline && <span>Due {formatDate(project.deadline)}</span>}
      </div>
    </Link>
  );
}
