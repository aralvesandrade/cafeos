import { useEffect, useState, useCallback } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { apiRequest } from '@/lib/api'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Table, TableHead, TableBody, TableRow, TableHeader, TableCell } from '@/components/ui/table'
import { ArrowLeft, Grid3X3, Ruler, HardHat, Map, User, Pencil, FileText, Trees } from 'lucide-react'

interface Producer {
  cpf: string
  name: string
  rg: string
  issuing_body: string
  gender: string
  birth_date: string | null
  marital_status: string
  phone: string
  email: string
  education: string
}

interface Farm {
  id: string
  name: string
  owner: string
  location: string
  total_area_ha: number
  planted_area_ha: number
  created_at: string
  phone?: string
  activities?: string
  main_crop?: string
  secondary_crop?: string
  state?: string
  city?: string
  address?: string
  production_system?: string
  commercialization_product?: string
  has_no_cnpj?: boolean
  cnpj?: string
  has_no_nirf?: boolean
  nirf?: string
  has_no_incra?: boolean
  incra_registration?: string
  has_no_state_registration?: boolean
  state_registration?: string
  has_no_dap?: boolean
  dap?: string
  has_no_car?: boolean
  car?: string
  fully_leased?: boolean
  land_value_per_ha?: number
  dam_area_ha?: number
  improvements_area_ha?: number
  roads_area_ha?: number
  app_area_ha?: number
  legal_reserve_area_ha?: number
  native_vegetation_area_ha?: number
  livestock_area_not_covered_ha?: number
  agriculture_area_not_covered_ha?: number
  non_agricultural_area_ha?: number
  producer?: Producer | null
}

interface Plot {
  id: string
  name: string
  area_ha: number
  cultivar: string
  soil_type: string
  altitude: number
  planting_year: number
}

export function FarmDetail() {
  const { farmId } = useParams<{ farmId: string }>()
  const navigate = useNavigate()
  const [farm, setFarm] = useState<Farm | null>(null)
  const [plots, setPlots] = useState<Plot[]>([])
  const [loading, setLoading] = useState(true)

  const load = useCallback(async () => {
    try {
      const [farmData, plotsData] = await Promise.all([
        apiRequest<Farm>(`/farms/${farmId}`),
        apiRequest<Plot[]>(`/farms/${farmId}/plots`),
      ])
      setFarm(farmData)
      setPlots(plotsData)
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }, [farmId])

  useEffect(() => { load() }, [load])

  if (loading) {
    return <div className="flex items-center justify-center h-64 text-muted-foreground">Carregando...</div>
  }

  if (!farm) {
    return (
      <div className="text-center py-16">
        <p className="text-muted-foreground mb-4">Fazenda não encontrada.</p>
        <Button variant="outline" onClick={() => navigate('/farms')}>Voltar</Button>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Button variant="ghost" size="sm" onClick={() => navigate('/farms')}>
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <h1 className="text-2xl font-bold text-primary">{farm.name}</h1>
        </div>
        <Button variant="outline" size="sm" asChild>
          <Link to={`/farms/${farmId}/edit`}>
            <Pencil className="h-4 w-4" />
            Editar
          </Link>
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <HardHat className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Proprietário</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">{farm.owner || '—'}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Map className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Localização</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">{farm.location || '—'}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Ruler className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Área Total</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">{farm.total_area_ha} ha</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <Ruler className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Área Plantada</CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-lg font-semibold text-foreground">{farm.planted_area_ha} ha</p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center gap-2 pb-2">
          <HardHat className="h-4 w-4 text-primary" />
          <CardTitle className="text-sm font-medium text-muted-foreground">Dados Gerais</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div><p className="text-muted-foreground">Telefone</p><p className="font-medium text-foreground">{farm.phone || '—'}</p></div>
            <div><p className="text-muted-foreground">Atividades</p><p className="font-medium text-foreground">{farm.activities || '—'}</p></div>
            <div><p className="text-muted-foreground">Cultura Principal</p><p className="font-medium text-foreground">{farm.main_crop || '—'}</p></div>
            <div><p className="text-muted-foreground">Cultura Secundária</p><p className="font-medium text-foreground">{farm.secondary_crop || '—'}</p></div>
            <div><p className="text-muted-foreground">Estado</p><p className="font-medium text-foreground">{farm.state || '—'}</p></div>
            <div><p className="text-muted-foreground">Cidade</p><p className="font-medium text-foreground">{farm.city || '—'}</p></div>
            <div><p className="text-muted-foreground">Endereço</p><p className="font-medium text-foreground">{farm.address || '—'}</p></div>
            <div><p className="text-muted-foreground">Sistema de Produção</p><p className="font-medium text-foreground">{farm.production_system || '—'}</p></div>
            <div><p className="text-muted-foreground">Produto de Comercialização</p><p className="font-medium text-foreground">{farm.commercialization_product || '—'}</p></div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center gap-2 pb-2">
          <FileText className="h-4 w-4 text-primary" />
          <CardTitle className="text-sm font-medium text-muted-foreground">Documentação</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4 text-sm">
            <div><p className="text-muted-foreground">CNPJ</p><p className="font-medium text-foreground">{farm.has_no_cnpj ? 'Não possui' : (farm.cnpj || '—')}</p></div>
            <div><p className="text-muted-foreground">NIRF</p><p className="font-medium text-foreground">{farm.has_no_nirf ? 'Não possui' : (farm.nirf || '—')}</p></div>
            <div><p className="text-muted-foreground">INCRA</p><p className="font-medium text-foreground">{farm.has_no_incra ? 'Não possui' : (farm.incra_registration || '—')}</p></div>
            <div><p className="text-muted-foreground">Inscrição Estadual</p><p className="font-medium text-foreground">{farm.has_no_state_registration ? 'Não possui' : (farm.state_registration || '—')}</p></div>
            <div><p className="text-muted-foreground">DAP</p><p className="font-medium text-foreground">{farm.has_no_dap ? 'Não possui' : (farm.dap || '—')}</p></div>
            <div><p className="text-muted-foreground">CAR</p><p className="font-medium text-foreground">{farm.has_no_car ? 'Não possui' : (farm.car || '—')}</p></div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className="flex flex-row items-center gap-2 pb-2">
          <Trees className="h-4 w-4 text-primary" />
          <CardTitle className="text-sm font-medium text-muted-foreground">Divisão de Área</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div><p className="text-muted-foreground">Totalmente Arrendada</p><p className="font-medium text-foreground">{farm.fully_leased ? 'Sim' : 'Não'}</p></div>
            <div><p className="text-muted-foreground">Valor da Terra/ha</p><p className="font-medium text-foreground">{farm.land_value_per_ha ? `R$ ${farm.land_value_per_ha.toFixed(2)}` : '—'}</p></div>
            <div><p className="text-muted-foreground">Área Plantada</p><p className="font-medium text-foreground">{farm.planted_area_ha ?? 0} ha</p></div>
            <div><p className="text-muted-foreground">Área Produtiva Disponível</p><p className="font-medium text-foreground">{(
              farm.total_area_ha
              - (farm.dam_area_ha || 0)
              - (farm.improvements_area_ha || 0)
              - (farm.roads_area_ha || 0)
              - (farm.app_area_ha || 0)
              - (farm.legal_reserve_area_ha || 0)
              - (farm.native_vegetation_area_ha || 0)
              - (farm.livestock_area_not_covered_ha || 0)
              - (farm.agriculture_area_not_covered_ha || 0)
              - (farm.non_agricultural_area_ha || 0)
            ).toFixed(2)} ha</p></div>
            <div><p className="text-muted-foreground">Açude</p><p className="font-medium text-foreground">{farm.dam_area_ha ?? 0} ha</p></div>
            <div><p className="text-muted-foreground">Benfeitorias</p><p className="font-medium text-foreground">{farm.improvements_area_ha ?? 0} ha</p></div>
            <div><p className="text-muted-foreground">Estradas</p><p className="font-medium text-foreground">{farm.roads_area_ha ?? 0} ha</p></div>
            <div><p className="text-muted-foreground">APP</p><p className="font-medium text-foreground">{farm.app_area_ha ?? 0} ha</p></div>
            <div><p className="text-muted-foreground">Reserva Legal</p><p className="font-medium text-foreground">{farm.legal_reserve_area_ha ?? 0} ha</p></div>
            <div><p className="text-muted-foreground">Vegetação Nativa</p><p className="font-medium text-foreground">{farm.native_vegetation_area_ha ?? 0} ha</p></div>
            <div><p className="text-muted-foreground">Pecuária Não Coberta</p><p className="font-medium text-foreground">{farm.livestock_area_not_covered_ha ?? 0} ha</p></div>
            <div><p className="text-muted-foreground">Agricultura Não Coberta</p><p className="font-medium text-foreground">{farm.agriculture_area_not_covered_ha ?? 0} ha</p></div>
            <div><p className="text-muted-foreground">Área Não Agrícola</p><p className="font-medium text-foreground">{farm.non_agricultural_area_ha ?? 0} ha</p></div>
          </div>
        </CardContent>
      </Card>

      {farm.producer && (
        <Card>
          <CardHeader className="flex flex-row items-center gap-2 pb-2">
            <User className="h-4 w-4 text-primary" />
            <CardTitle className="text-sm font-medium text-muted-foreground">Produtor</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
              <div>
                <p className="text-muted-foreground">Nome</p>
                <p className="font-medium text-foreground">{farm.producer.name || '—'}</p>
              </div>
              <div>
                <p className="text-muted-foreground">CPF</p>
                <p className="font-medium text-foreground">{farm.producer.cpf || '—'}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Telefone</p>
                <p className="font-medium text-foreground">{farm.producer.phone || '—'}</p>
              </div>
              <div>
                <p className="text-muted-foreground">E-mail</p>
                <p className="font-medium text-foreground">{farm.producer.email || '—'}</p>
              </div>
              <div>
                <p className="text-muted-foreground">RG</p>
                <p className="font-medium text-foreground">{farm.producer.rg || '—'}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Órgão Emissor</p>
                <p className="font-medium text-foreground">{farm.producer.issuing_body || '—'}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Sexo</p>
                <p className="font-medium text-foreground">{farm.producer.gender || '—'}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Data de Nascimento</p>
                <p className="font-medium text-foreground">{farm.producer.birth_date ? new Date(farm.producer.birth_date).toLocaleDateString('pt-BR') : '—'}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Estado Civil</p>
                <p className="font-medium text-foreground">{farm.producer.marital_status || '—'}</p>
              </div>
              <div>
                <p className="text-muted-foreground">Escolaridade</p>
                <p className="font-medium text-foreground">{farm.producer.education || '—'}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      <div>
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-primary flex items-center gap-2">
            <Grid3X3 className="h-5 w-5" />
            Talhões
          </h2>
          <Button variant="outline" size="sm" onClick={() => navigate('/plots')}>
            Gerenciar Talhões
          </Button>
        </div>

        <Table>
          <TableHead>
            <TableRow>
              <TableHeader>Nome</TableHeader>
              <TableHeader>Área</TableHeader>
              <TableHeader>Cultivar</TableHeader>
              <TableHeader>Solo</TableHeader>
              <TableHeader>Altitude</TableHeader>
              <TableHeader>Ano Plantio</TableHeader>
              <TableHeader className="text-right">Ações</TableHeader>
            </TableRow>
          </TableHead>
          <TableBody>
            {plots.map((plot) => (
              <TableRow key={plot.id}>
                <TableCell className="font-medium">{plot.name}</TableCell>
                <TableCell>{plot.area_ha} ha</TableCell>
                <TableCell>{plot.cultivar}</TableCell>
                <TableCell>{plot.soil_type}</TableCell>
                <TableCell>{plot.altitude} m</TableCell>
                <TableCell>{plot.planting_year}</TableCell>
                <TableCell className="text-right">
                  <Button variant="ghost" size="sm" asChild>
                    <Link to={`/plots/${plot.id}`}>Detalhes</Link>
                  </Button>
                </TableCell>
              </TableRow>
            ))}
            {plots.length === 0 && (
              <TableRow>
                <TableCell colSpan={7} className="text-center text-muted-foreground py-8">
                  Nenhum talhão cadastrado nesta fazenda.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
