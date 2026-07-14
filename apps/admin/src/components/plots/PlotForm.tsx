import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'

export interface PlotData {
  name: string
  farm_id: string
  area: string
  cultivar: string
  soil_type: string
  altitude: string
  planting_year: string

  leased: boolean
  stage: string
  irrigation: string

  activation_date: string
  planting_date: string
  deactivation_date: string

  intercropped: boolean
  secondary_crop: string
  notes: string

  crop_type: string

  formation_cost_per_ha: string
  useful_life_years: string
  row_spacing_m: string
  plant_spacing_m: string
  plant_count: string

  dam_area: string
  improvements_area: string
  roads_area: string
  app_area: string
  legal_reserve_area: string
}

export const emptyPlot: PlotData = {
  name: '',
  farm_id: '',
  area: '',
  cultivar: '',
  soil_type: '',
  altitude: '',
  planting_year: '',
  leased: false,
  stage: 'formacao',
  irrigation: '',
  activation_date: '',
  planting_date: '',
  deactivation_date: '',
  intercropped: false,
  secondary_crop: '',
  notes: '',
  crop_type: '',
  formation_cost_per_ha: '',
  useful_life_years: '',
  row_spacing_m: '',
  plant_spacing_m: '',
  plant_count: '',
  dam_area: '',
  improvements_area: '',
  roads_area: '',
  app_area: '',
  legal_reserve_area: '',
}

interface Farm { id: string; name: string }

interface PlotFormProps {
  initial?: PlotData
  farms: Farm[]
  onSave: (data: PlotData) => void
  onCancel: () => void
  loading?: boolean
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="space-y-3">
      <h3 className="text-sm font-semibold text-primary border-b border-border pb-1">{title}</h3>
      {children}
    </div>
  )
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div>
      <label className="block text-sm font-medium text-foreground mb-1">{label}</label>
      {children}
    </div>
  )
}

export function PlotForm({ initial, farms, onSave, onCancel, loading }: PlotFormProps) {
  const [form, setForm] = useState<PlotData>(emptyPlot)

  useEffect(() => {
    setForm(initial ? { ...emptyPlot, ...initial } : emptyPlot)
  }, [initial])

  const set = <K extends keyof PlotData>(key: K, value: PlotData[K]) => setForm((f) => ({ ...f, [key]: value }))

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSave(form)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <Section title="Informações do talhão">
        <div className="grid grid-cols-2 gap-4">
          <label className="flex items-center gap-2 text-sm text-foreground col-span-2">
            <input type="checkbox" checked={form.leased} onChange={(e) => set('leased', e.target.checked)} />
            Arrendado?
          </label>
          <Field label="Nome do talhão">
            <Input value={form.name} onChange={(e) => set('name', e.target.value)} required />
          </Field>
          <Field label="Fazenda">
            <Select value={form.farm_id} onChange={(e) => set('farm_id', e.target.value)} required>
              <option value="">Selecione...</option>
              {farms.map((f) => <option key={f.id} value={f.id}>{f.name}</option>)}
            </Select>
          </Field>
          <Field label="Área (ha)">
            <Input type="number" step="0.01" value={form.area} onChange={(e) => set('area', e.target.value)} required />
          </Field>
          <Field label="Estágio">
            <Select value={form.stage} onChange={(e) => set('stage', e.target.value)} required>
              <option value="formacao">Formação</option>
              <option value="producao">Produção</option>
            </Select>
          </Field>
          <Field label="Irrigação">
            <Input value={form.irrigation} onChange={(e) => set('irrigation', e.target.value)} placeholder="Ex: Não irrigado" />
          </Field>
          <Field label="Data Ativação">
            <Input type="date" value={form.activation_date} onChange={(e) => set('activation_date', e.target.value)} />
          </Field>
          <Field label="Data de Plantio">
            <Input type="date" value={form.planting_date} onChange={(e) => set('planting_date', e.target.value)} />
          </Field>
          <Field label="Data Desativação">
            <Input type="date" value={form.deactivation_date} onChange={(e) => set('deactivation_date', e.target.value)} />
          </Field>
          <label className="flex items-center gap-2 text-sm text-foreground">
            <input type="checkbox" checked={form.intercropped} onChange={(e) => set('intercropped', e.target.checked)} />
            Consorciada
          </label>
          <Field label="Cultura Secundária">
            <Input value={form.secondary_crop} onChange={(e) => set('secondary_crop', e.target.value)} />
          </Field>
          <Field label="Ano Plantio">
            <Input type="number" value={form.planting_year} onChange={(e) => set('planting_year', e.target.value)} />
          </Field>
          <Field label="Altitude (m)">
            <Input type="number" value={form.altitude} onChange={(e) => set('altitude', e.target.value)} />
          </Field>
          <Field label="Tipo de Solo">
            <Select value={form.soil_type} onChange={(e) => set('soil_type', e.target.value)}>
              <option value="">Selecione...</option>
              <option value="argiloso">Argiloso</option>
              <option value="arenoso">Arenoso</option>
              <option value="siltoso">Siltoso</option>
              <option value="organico">Orgânico</option>
            </Select>
          </Field>
          <div className="col-span-2">
            <Field label="Observações">
              <Input value={form.notes} onChange={(e) => set('notes', e.target.value)} />
            </Field>
          </div>
        </div>
      </Section>

      <Section title="Informações da atividade">
        <div className="grid grid-cols-3 gap-4">
          <Field label="Tipo">
            <Select value={form.crop_type} onChange={(e) => set('crop_type', e.target.value)}>
              <option value="">Selecione...</option>
              <option value="arabica">Arábica</option>
              <option value="robusta">Robusta</option>
            </Select>
          </Field>
          <Field label="Variedade / Clone">
            <Input value={form.cultivar} onChange={(e) => set('cultivar', e.target.value)} />
          </Field>
          <Field label="Vida útil (anos)">
            <Input type="number" value={form.useful_life_years} onChange={(e) => set('useful_life_years', e.target.value)} />
          </Field>
          <Field label="Custo de formação/ha (R$)">
            <Input type="number" step="0.01" value={form.formation_cost_per_ha} onChange={(e) => set('formation_cost_per_ha', e.target.value)} />
          </Field>
          <Field label="Espaçamento Linha (m)">
            <Input type="number" step="0.01" value={form.row_spacing_m} onChange={(e) => set('row_spacing_m', e.target.value)} />
          </Field>
          <Field label="Espaçamento Planta (m)">
            <Input type="number" step="0.01" value={form.plant_spacing_m} onChange={(e) => set('plant_spacing_m', e.target.value)} />
          </Field>
          <Field label="Nº de plantas">
            <Input type="number" value={form.plant_count} onChange={(e) => set('plant_count', e.target.value)} />
          </Field>
        </div>
      </Section>

      <Section title="Divisão de área do talhão">
        <div className="grid grid-cols-3 gap-4">
          <Field label="Açude ou Represa (ha)">
            <Input type="number" step="0.01" value={form.dam_area} onChange={(e) => set('dam_area', e.target.value)} />
          </Field>
          <Field label="Benfeitorias (ha)">
            <Input type="number" step="0.01" value={form.improvements_area} onChange={(e) => set('improvements_area', e.target.value)} />
          </Field>
          <Field label="Estradas (ha)">
            <Input type="number" step="0.01" value={form.roads_area} onChange={(e) => set('roads_area', e.target.value)} />
          </Field>
          <Field label="APP (ha)">
            <Input type="number" step="0.01" value={form.app_area} onChange={(e) => set('app_area', e.target.value)} />
          </Field>
          <Field label="Reserva legal (ha)">
            <Input type="number" step="0.01" value={form.legal_reserve_area} onChange={(e) => set('legal_reserve_area', e.target.value)} />
          </Field>
        </div>
      </Section>

      <div className="flex justify-end gap-3 pt-2 sticky bottom-0 bg-background pb-2">
        <Button type="button" variant="outline" onClick={onCancel}>Cancelar</Button>
        <Button type="submit" disabled={loading}>
          {loading ? 'Salvando...' : 'Salvar'}
        </Button>
      </div>
    </form>
  )
}
