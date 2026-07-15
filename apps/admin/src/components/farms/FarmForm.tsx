import { useState, useEffect } from 'react'
import { apiRequest } from '@/lib/api'
import { useRoles } from '@/lib/roles'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import { Plus, Trash2 } from 'lucide-react'

export interface ProducerData {
  user_id: string
  role_id: string
  cpf: string
  name: string
  rg: string
  issuing_body: string
  gender: string
  birth_date: string
  marital_status: string
  phone: string
  email: string
  education: string
}

interface OrgUser {
  id: string
  name: string
  email: string
}

export interface FarmData {
  name: string
  owner: string
  location: string
  total_area: string
  planted_area: string

  phone: string
  activities: string
  main_crop: string
  secondary_crop: string
  state: string
  city: string
  address: string
  production_system: string
  commercialization_product: string

  has_no_cnpj: boolean
  cnpj: string
  has_no_nirf: boolean
  nirf: string
  has_no_incra: boolean
  incra_registration: string
  has_no_state_registration: boolean
  state_registration: string
  has_no_dap: boolean
  dap: string
  has_no_car: boolean
  car: string

  fully_leased: boolean
  land_value_per_ha: string

  dam_area: string
  improvements_area: string
  roads_area: string
  app_area: string
  legal_reserve_area: string
  native_vegetation_area: string
  livestock_area_not_covered: string
  agriculture_area_not_covered: string
  non_agricultural_area: string

  producers: ProducerData[]
}

export const emptyProducer: ProducerData = {
  user_id: '',
  role_id: '',
  cpf: '',
  name: '',
  rg: '',
  issuing_body: '',
  gender: '',
  birth_date: '',
  marital_status: '',
  phone: '',
  email: '',
  education: '',
}

export const emptyFarm: FarmData = {
  name: '',
  owner: '',
  location: '',
  total_area: '',
  planted_area: '',
  phone: '',
  activities: '',
  main_crop: '',
  secondary_crop: '',
  state: '',
  city: '',
  address: '',
  production_system: '',
  commercialization_product: '',
  has_no_cnpj: false,
  cnpj: '',
  has_no_nirf: false,
  nirf: '',
  has_no_incra: false,
  incra_registration: '',
  has_no_state_registration: false,
  state_registration: '',
  has_no_dap: false,
  dap: '',
  has_no_car: false,
  car: '',
  fully_leased: false,
  land_value_per_ha: '',
  dam_area: '',
  improvements_area: '',
  roads_area: '',
  app_area: '',
  legal_reserve_area: '',
  native_vegetation_area: '',
  livestock_area_not_covered: '',
  agriculture_area_not_covered: '',
  non_agricultural_area: '',
  producers: [],
}

interface FarmFormProps {
  initial?: FarmData
  onSave: (data: FarmData) => void
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

function DocField({
  label,
  hasNone,
  value,
  onHasNoneChange,
  onValueChange,
}: {
  label: string
  hasNone: boolean
  value: string
  onHasNoneChange: (v: boolean) => void
  onValueChange: (v: string) => void
}) {
  return (
    <div>
      <div className="flex items-center justify-between mb-1">
        <label className="block text-sm font-medium text-foreground">{label}</label>
        <label className="flex items-center gap-1 text-xs text-muted-foreground">
          <input type="checkbox" checked={hasNone} onChange={(e) => onHasNoneChange(e.target.checked)} />
          Não possui
        </label>
      </div>
      <Input value={value} disabled={hasNone} onChange={(e) => onValueChange(e.target.value)} />
    </div>
  )
}

export function FarmForm({ initial, onSave, onCancel, loading }: FarmFormProps) {
  const [form, setForm] = useState<FarmData>(emptyFarm)
  const [orgUsers, setOrgUsers] = useState<OrgUser[]>([])
  const { roles } = useRoles()

  useEffect(() => {
    apiRequest<OrgUser[]>('/users').then(setOrgUsers).catch(console.error)
  }, [])

  useEffect(() => {
    if (initial) setForm({ ...emptyFarm, ...initial, producers: initial.producers ?? [] })
  }, [initial])

  const set = <K extends keyof FarmData>(key: K, value: FarmData[K]) => setForm((f) => ({ ...f, [key]: value }))

  const addProducer = () => setForm((f) => ({ ...f, producers: [...f.producers, { ...emptyProducer }] }))
  const removeProducer = (index: number) =>
    setForm((f) => ({ ...f, producers: f.producers.filter((_, i) => i !== index) }))
  const setProducer = <K extends keyof ProducerData>(index: number, key: K, value: ProducerData[K]) =>
    setForm((f) => ({
      ...f,
      producers: f.producers.map((p, i) => (i === index ? { ...p, [key]: value } : p)),
    }))

  const roleKeyFor = (roleId: string) => roles.find((r) => r.id === roleId)?.key

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    onSave(form)
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <Section title="Propriedade">
        <div className="grid grid-cols-2 gap-4">
          <Field label="Nome da Propriedade">
            <Input value={form.name} onChange={(e) => set('name', e.target.value)} required />
          </Field>
          <Field label="Telefone">
            <Input value={form.phone} onChange={(e) => set('phone', e.target.value)} />
          </Field>
          <Field label="Proprietário">
            <Input value={form.owner} onChange={(e) => set('owner', e.target.value)} required />
          </Field>
          <Field label="Atividades">
            <Input value={form.activities} onChange={(e) => set('activities', e.target.value)} placeholder="Ex: CAFEICULTURA" />
          </Field>
          <Field label="Cultura Principal">
            <Input value={form.main_crop} onChange={(e) => set('main_crop', e.target.value)} />
          </Field>
          <Field label="Cultura Secundária">
            <Input value={form.secondary_crop} onChange={(e) => set('secondary_crop', e.target.value)} />
          </Field>
          <Field label="UF">
            <Input value={form.state} onChange={(e) => set('state', e.target.value)} maxLength={2} />
          </Field>
          <Field label="Município">
            <Input value={form.city} onChange={(e) => set('city', e.target.value)} />
          </Field>
          <Field label="Endereço">
            <Input value={form.address} onChange={(e) => set('address', e.target.value)} />
          </Field>
          <Field label="Localização (referência)">
            <Input value={form.location} onChange={(e) => set('location', e.target.value)} required />
          </Field>
          <Field label="Sistema de Produção">
            <Input value={form.production_system} onChange={(e) => set('production_system', e.target.value)} />
          </Field>
          <Field label="Produto de Comercialização">
            <Input value={form.commercialization_product} onChange={(e) => set('commercialization_product', e.target.value)} />
          </Field>
        </div>
      </Section>

      <Section title="Documentação">
        <div className="grid grid-cols-2 gap-4">
          <DocField label="CNPJ" hasNone={form.has_no_cnpj} value={form.cnpj} onHasNoneChange={(v) => set('has_no_cnpj', v)} onValueChange={(v) => set('cnpj', v)} />
          <DocField label="Número do NIRF" hasNone={form.has_no_nirf} value={form.nirf} onHasNoneChange={(v) => set('has_no_nirf', v)} onValueChange={(v) => set('nirf', v)} />
          <DocField label="Inscrição no Incra" hasNone={form.has_no_incra} value={form.incra_registration} onHasNoneChange={(v) => set('has_no_incra', v)} onValueChange={(v) => set('incra_registration', v)} />
          <DocField label="Inscrição Estadual" hasNone={form.has_no_state_registration} value={form.state_registration} onHasNoneChange={(v) => set('has_no_state_registration', v)} onValueChange={(v) => set('state_registration', v)} />
          <DocField label="DAP" hasNone={form.has_no_dap} value={form.dap} onHasNoneChange={(v) => set('has_no_dap', v)} onValueChange={(v) => set('dap', v)} />
          <DocField label="CAR" hasNone={form.has_no_car} value={form.car} onHasNoneChange={(v) => set('has_no_car', v)} onValueChange={(v) => set('car', v)} />
        </div>
      </Section>

      <Section title="Dados gerais da propriedade">
        <div className="grid grid-cols-3 gap-4 items-end">
          <label className="flex items-center gap-2 text-sm text-foreground">
            <input type="checkbox" checked={form.fully_leased} onChange={(e) => set('fully_leased', e.target.checked)} />
            Propriedade totalmente arrendada?
          </label>
          <Field label="Valor da terra nua/ha (R$)">
            <Input type="number" step="0.01" value={form.land_value_per_ha} onChange={(e) => set('land_value_per_ha', e.target.value)} />
          </Field>
          <Field label="Área Total da Propriedade (ha)">
            <Input type="number" step="0.01" value={form.total_area} onChange={(e) => set('total_area', e.target.value)} required />
          </Field>
        </div>
      </Section>

      <Section title="Divisão das áreas da propriedade">
        <div className="grid grid-cols-3 gap-4">
          <Field label="Área Plantada (ha)">
            <Input type="number" step="0.01" value={form.planted_area} onChange={(e) => set('planted_area', e.target.value)} required />
          </Field>
          <Field label="Açude ou represa (ha)">
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
          <Field label="Reserva Legal (ha)">
            <Input type="number" step="0.01" value={form.legal_reserve_area} onChange={(e) => set('legal_reserve_area', e.target.value)} />
          </Field>
          <Field label="Vegetação Nativa (ha)">
            <Input type="number" step="0.01" value={form.native_vegetation_area} onChange={(e) => set('native_vegetation_area', e.target.value)} />
          </Field>
          <Field label="Área Pecuária não contemplada (ha)">
            <Input type="number" step="0.01" value={form.livestock_area_not_covered} onChange={(e) => set('livestock_area_not_covered', e.target.value)} />
          </Field>
          <Field label="Área Agricultura não contemplada (ha)">
            <Input type="number" step="0.01" value={form.agriculture_area_not_covered} onChange={(e) => set('agriculture_area_not_covered', e.target.value)} />
          </Field>
          <Field label="Áreas não Agrícolas (ha)">
            <Input type="number" step="0.01" value={form.non_agricultural_area} onChange={(e) => set('non_agricultural_area', e.target.value)} />
          </Field>
        </div>
      </Section>

      <Section title="Vínculos da Fazenda">
        <div className="space-y-4">
          {form.producers.map((producer, index) => {
            const isProprietario = roleKeyFor(producer.role_id) === 'proprietario'
            return (
              <div key={index} className="rounded-lg border border-border p-4 space-y-4">
                <div className="flex items-center justify-between">
                  <div className="grid grid-cols-2 gap-4 flex-1">
                    <Field label="Usuário">
                      <Select value={producer.user_id} onChange={(e) => setProducer(index, 'user_id', e.target.value)} required>
                        <option value="">Selecione...</option>
                        {orgUsers.map((u) => (
                          <option key={u.id} value={u.id}>{u.name} ({u.email})</option>
                        ))}
                      </Select>
                    </Field>
                    <Field label="Papel nesta fazenda">
                      <Select value={producer.role_id} onChange={(e) => setProducer(index, 'role_id', e.target.value)} required>
                        <option value="">Selecione...</option>
                        {roles.map((r) => (
                          <option key={r.id} value={r.id}>{r.name}</option>
                        ))}
                      </Select>
                    </Field>
                  </div>
                  <Button type="button" variant="ghost" size="sm" className="ml-4 mt-6" onClick={() => removeProducer(index)}>
                    <Trash2 className="h-4 w-4 text-destructive" />
                  </Button>
                </div>

                {isProprietario && (
                  <div className="grid grid-cols-3 gap-4 pt-2 border-t border-border">
                    <Field label="CPF">
                      <Input value={producer.cpf} onChange={(e) => setProducer(index, 'cpf', e.target.value)} />
                    </Field>
                    <Field label="Nome">
                      <Input value={producer.name} onChange={(e) => setProducer(index, 'name', e.target.value)} />
                    </Field>
                    <Field label="RG">
                      <Input value={producer.rg} onChange={(e) => setProducer(index, 'rg', e.target.value)} />
                    </Field>
                    <Field label="Órgão Emissor">
                      <Input value={producer.issuing_body} onChange={(e) => setProducer(index, 'issuing_body', e.target.value)} />
                    </Field>
                    <Field label="Sexo">
                      <Select value={producer.gender} onChange={(e) => setProducer(index, 'gender', e.target.value)}>
                        <option value="">Selecione...</option>
                        <option value="masculino">Masculino</option>
                        <option value="feminino">Feminino</option>
                      </Select>
                    </Field>
                    <Field label="Data de nascimento">
                      <Input type="date" value={producer.birth_date} onChange={(e) => setProducer(index, 'birth_date', e.target.value)} />
                    </Field>
                    <Field label="Estado civil">
                      <Select value={producer.marital_status} onChange={(e) => setProducer(index, 'marital_status', e.target.value)}>
                        <option value="">Selecione...</option>
                        <option value="solteiro">Solteiro(a)</option>
                        <option value="casado">Casado(a)</option>
                        <option value="divorciado">Divorciado(a)</option>
                        <option value="viuvo">Viúvo(a)</option>
                      </Select>
                    </Field>
                    <Field label="Telefone">
                      <Input value={producer.phone} onChange={(e) => setProducer(index, 'phone', e.target.value)} />
                    </Field>
                    <Field label="E-mail">
                      <Input type="email" value={producer.email} onChange={(e) => setProducer(index, 'email', e.target.value)} />
                    </Field>
                    <Field label="Escolaridade">
                      <Select value={producer.education} onChange={(e) => setProducer(index, 'education', e.target.value)}>
                        <option value="">Selecione...</option>
                        <option value="fundamental">Ensino Fundamental</option>
                        <option value="medio">Ensino Médio</option>
                        <option value="superior">Ensino Superior</option>
                      </Select>
                    </Field>
                  </div>
                )}
              </div>
            )
          })}
          <Button type="button" variant="outline" onClick={addProducer}>
            <Plus className="h-4 w-4" /> Adicionar vínculo
          </Button>
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
