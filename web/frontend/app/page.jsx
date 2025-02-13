"use client"
import { useEffect, useState } from "react";
import { ToastContainer, toast } from "react-toastify";

export default function Page() {
  const [user, setUser] = useState('')
  const [creci, setCreci] = useState('')
  const [content, setContent] = useState('')
  const [parecerId, setParecerId] = useState('')

  const [date, setDate] = useState('')
  const [editParecer, setEditParecer] = useState(false)
  const [pareceres, setPareceres] = useState([])
  const [isLoading, setIsLoading] = useState(true)
  const [reload, setReload] = useState(false)


  const notify = (error) => {
    toast(error)
  }

  const fillEditParecer = (parecer) => {
    setUser(parecer.user)
    setCreci(parecer.creci)
    setDate(new Date(parecer.date).toISOString().split('T')[0])
    setContent(parecer.content)
    setParecerId(parecer.id)
    setEditParecer(true)
  }

  const createParecer = () => {
    fetch(`${process.env.NEXT_PUBLIC_BACKEND_API_URL}/parecer`, {
      method: 'POST',
      body: JSON.stringify({
        user, creci, content, date
      })
    })
      .then(res => {
        if (res.status != 201) {
          notify("Erro no para criar parecer, codigo " + res.status)
        }
        setUser('')
        setCreci('')
        setDate('')
        setContent('')
        setParecerId('')
        setIsLoading(true)
        setReload(!reload)
      })
      .catch(err => notify('Erro ao mandar parecer ' + err))
  }

  const getAllPareceres = () => {
    fetch(`${process.env.NEXT_PUBLIC_BACKEND_API_URL}/parecer`)
      .then(res => {
        if (res.status != 200) {
          notify("Erro no servidor para carregar pareceres")
        }
        res.json()
          .then(data => {
            setPareceres(data)
            setIsLoading(false)
          })
          .catch(err => notify('Erro ao carregar pareceres ' + err))
      })
      .catch(err => {
        notify('Erro ao carregar pareceres ' + err)
      })
  }

  const updateParecer = () => {
    fetch(`${process.env.NEXT_PUBLIC_BACKEND_API_URL}/parecer?id=${parecerId}`, {
      method: 'PUT',
      body: JSON.stringify({
        user, creci, date, content
      })
    }).then(res => {
      if (res.status != 200) {
        notify("Erro no servidor para atualizar parecer")
      }
      setUser('')
      setCreci('')
      setContent('')
      setDate('')
      setParecerId('')
      setIsLoading(true)
      setReload(!reload)
    }).catch(err => {
      notify('Erro ao editar parecer ' + err)
      setUser('')
      setCreci('')
      setDate('')
      setContent('')
      setParecerId('')
      setIsLoading(true)
      setEditParecer(false)
    })
  }

  const deleteParecer = (parecer) => {
    fetch(`${process.env.NEXT_PUBLIC_BACKEND_API_URL}/parecer?id=${parecer.id}`, {
      method: 'DELETE',
    }).then(res => {
      if (res.status != 200) {
        notify("Erro no servidor para deletar parecer")
      }
      setIsLoading(true)
      setReload(!reload)
    }).catch(err => notify('Erro ao deletar parecer ' + err))
  }


  useEffect(() => {
    getAllPareceres()
  }, [reload])

  return (
    <>
      <ToastContainer aria-label={undefined} />

      <h1 className="text-3xl">Gerador de Parecer</h1>

      <div className="flex flex-col justify-center w-full px-10 sm:w-1/2 gap-4">
        <div className="flex flex-col">
          <label>Usuário:</label>
          <input className="rounded p-2 dark:bg-gray-700 bg-gray-100" type="text" name="user" id="user" value={user} onChange={(e) => setUser(e.target.value)} />
        </div>

        <div className="flex flex-col">
          <label>Creci:</label>
          <input className="rounded p-2 dark:bg-gray-700 bg-gray-100" type="text" name="creci" id="creci" value={creci} onChange={(e) => setCreci(e.target.value)} />
        </div>

        <div className="flex flex-col">
          <label>Date da emissão:</label>
          <input className="rounded p-2 dark:bg-gray-700 bg-gray-100" type="date" name="creci" id="creci" value={date} onChange={(e) => setDate(e.target.value)} />
        </div>

        <div className="flex flex-col">
          <label>Conteúdo:</label>
          <textarea className="rounded p-2 dark:bg-gray-700 bg-gray-100" name="content" id="content" rows={10} value={content} onChange={(e) => setContent(e.target.value)}></textarea>
        </div>

        <div className="flex justify-center items-center my-4">
          {!editParecer ? <button onClick={() => {
            if (user == '' || creci == '' || content == '') {
              notify('Preencha todos os campos')
            } else {
              createParecer()
            }
          }} className="flex border rounded-xl dark:bg-gray-700 bg-gray-100 py-2 px-4">Gerar Parecer</button> :
            <button onClick={updateParecer} className="flex border rounded-xl dark:bg-gray-700 bg-gray-100 py-2 px-4">Editar Parecer</button>
          }
        </div>

        <ul className="flex justify-center flex-col gap-y-4">
          {isLoading ? <p className="text-center">Carregando...</p> : pareceres &&
            pareceres.map((parecer) => (
              <li key={parecer.id} className="flex gap-4  w-full">
                <div className="p-4 border rounded bg-gray-100 w-full dark:bg-gray-700">
                  <p><strong>Usuário:</strong> {parecer.user}</p>
                  <p><strong>Creci:</strong> {parecer.creci}</p>
                  <p><strong>Data:</strong> {parecer.date}</p>
                </div>
                <div className="flex flex-col md:flex-row gap-4 justify-evenly">
                  <a className="flex items-center p-4 border rounded bg-gray-100 dark:bg-gray-700 " href={`${process.env.NEXT_PUBLIC_BACKEND_API_URL}/parecer?id=${parecer.id}`} target="_blank">Baixar</a>
                  <div className="flex items-center p-4 border rounded bg-gray-100 dark:bg-gray-700" onClick={() => {
                    fillEditParecer(parecer)
                    setEditParecer(true)
                    window.scrollTo(0, 0)
                  }}>Editar</div>
                  <div className="flex items-center p-4 border rounded bg-gray-100 dark:bg-gray-700" onClick={() => {
                    deleteParecer(parecer)
                  }}>Deletar</div>
                </div>
              </li>
            ))}
        </ul>
      </div >

    </>
  );
}