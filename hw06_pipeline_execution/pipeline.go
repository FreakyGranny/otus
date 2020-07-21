package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

// Stage ...
type Stage func(in In) (out Out)

func wrapStage(done In, in In, fn Stage) Out {
	outStream := make(Bi)
	go func() {
		defer close(outStream)
		for {
			select {
			case <-done:
				return
			case x, open := <-in:
				if open {
					outStream <- x
				} else {
					return
				}
			}
		}
	}()
	return fn(outStream)
}

// ExecutePipeline execute given pipeline.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stream := wrapStage(done, in, func(in In) Out { return in })

	for _, v := range stages {
		stream = wrapStage(done, stream, v)
	}

	return stream
}
