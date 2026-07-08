import { Link } from 'react-router-dom'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Pencil, Trash2, ExternalLink } from 'lucide-react'
import type { Farm } from '@/pages/Farms'

interface FarmListProps {
  farms: Farm[]
  onEdit: (farm: Farm) => void
  onDelete: (id: string) => void
}

export function FarmList({ farms, onEdit, onDelete }: FarmListProps) {
  return (
    <Table>
      <TableHead>
        <TableRow>
          <TableHeader>Nome</TableHeader>
          <TableHeader>Proprietário</TableHeader>
          <TableHeader>Área Total</TableHeader>
          <TableHeader>Área Plantada</TableHeader>
          <TableHeader>Localização</TableHeader>
          <TableHeader className="text-right">Ações</TableHeader>
        </TableRow>
      </TableHead>
      <TableBody>
        {farms.map((farm) => (
          <TableRow key={farm.id}>
            <TableCell className="font-medium">
              <Link to={`/farms/${farm.id}`} className="text-coffee-green hover:underline">
                {farm.name}
              </Link>
            </TableCell>
            <TableCell>{farm.owner}</TableCell>
            <TableCell>{farm.total_area_ha} ha</TableCell>
            <TableCell>{farm.planted_area_ha} ha</TableCell>
            <TableCell>{farm.location}</TableCell>
            <TableCell className="text-right">
              <div className="flex justify-end gap-1">
                <Button variant="ghost" size="sm" asChild>
                  <Link to={`/farms/${farm.id}`}>
                    <ExternalLink className="h-4 w-4" />
                  </Link>
                </Button>
                <Button variant="ghost" size="sm" onClick={() => onEdit(farm)}>
                  <Pencil className="h-4 w-4" />
                </Button>
                <Button variant="ghost" size="sm" onClick={() => onDelete(farm.id)}>
                  <Trash2 className="h-4 w-4 text-red-500" />
                </Button>
              </div>
            </TableCell>
          </TableRow>
        ))}
        {farms.length === 0 && (
          <TableRow>
            <TableCell colSpan={6} className="text-center text-coffee-text-light py-8">
              Nenhuma fazenda cadastrada.
            </TableCell>
          </TableRow>
        )}
      </TableBody>
    </Table>
  )
}
